package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB wraps sql.DB to keep the health-check interface context-aware.
type DB struct {
	*sql.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func Open(ctx context.Context, databaseURL string) (*DB, error) {
	dsn, err := mysqlDSN(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	pool.SetConnMaxLifetime(3 * time.Minute)
	pool.SetMaxIdleConns(5)
	pool.SetMaxOpenConns(10)

	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &DB{DB: pool}, nil
}

// mysqlDSN converts the mysql:// URL supplied by Aiven into the native DSN
// expected by go-sql-driver/mysql. Aiven's ssl-mode=REQUIRED is mapped to a
// verified TLS connection; ssl-mode=DISABLED is intended for local development.
func mysqlDSN(databaseURL string) (string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}
	if u.Scheme != "mysql" {
		return "", fmt.Errorf("unsupported scheme %q; expected mysql", u.Scheme)
	}
	if u.User == nil || u.User.Username() == "" {
		return "", fmt.Errorf("database username is required")
	}
	password, ok := u.User.Password()
	if !ok {
		return "", fmt.Errorf("database password is required")
	}
	if u.Hostname() == "" || u.Port() == "" {
		return "", fmt.Errorf("database host and port are required")
	}
	databaseName := strings.TrimPrefix(u.EscapedPath(), "/")
	if databaseName == "" {
		return "", fmt.Errorf("database name is required")
	}
	databaseName, err = url.PathUnescape(databaseName)
	if err != nil {
		return "", fmt.Errorf("decode database name: %w", err)
	}

	sslMode := strings.ToUpper(u.Query().Get("ssl-mode"))
	var tlsConfig string
	switch sslMode {
	case "REQUIRED", "VERIFY_CA", "VERIFY_IDENTITY":
		tlsConfig = "true"
	case "DISABLED":
		tlsConfig = "false"
	default:
		return "", fmt.Errorf("ssl-mode must be REQUIRED or DISABLED")
	}

	cfg := mysql.NewConfig()
	cfg.User = u.User.Username()
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(u.Hostname(), u.Port())
	cfg.DBName = databaseName
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.TLSConfig = tlsConfig
	cfg.Timeout = 10 * time.Second
	cfg.ReadTimeout = 10 * time.Second
	cfg.WriteTimeout = 10 * time.Second
	return cfg.FormatDSN(), nil
}
