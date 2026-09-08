package database

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB wraps sql.DB to keep the health-check interface context-aware.
type DB struct {
	*sql.DB
}

var (
	tlsConfigMu         sync.Mutex
	registeredTLSConfig = make(map[string]struct{})
)

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func Open(ctx context.Context, databaseURL, caCertFile string) (*DB, error) {
	dsn, err := mysqlDSN(databaseURL, caCertFile)
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
// expected by go-sql-driver/mysql. ssl-mode=REQUIRED encrypts without verifying
// the issuer, matching MySQL's mode semantics. Supplying a CA file enables full
// certificate and hostname verification.
func mysqlDSN(databaseURL, caCertFile string) (string, error) {
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
	case "REQUIRED":
		tlsConfig = "skip-verify"
	case "VERIFY_CA", "VERIFY_IDENTITY":
		tlsConfig = "true"
	case "DISABLED":
		tlsConfig = "false"
	default:
		return "", fmt.Errorf("unsupported ssl-mode %q", sslMode)
	}
	if caCertFile != "" && sslMode != "DISABLED" {
		tlsConfig, err = registerTLSConfig(caCertFile, u.Hostname())
		if err != nil {
			return "", err
		}
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

func registerTLSConfig(caCertFile, serverName string) (string, error) {
	pem, err := os.ReadFile(caCertFile)
	if err != nil {
		return "", fmt.Errorf("read database CA certificate: %w", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pem) {
		return "", fmt.Errorf("parse database CA certificate: no certificates found")
	}

	digest := sha256.Sum256(append(pem, serverName...))
	name := fmt.Sprintf("aiven-%x", digest[:8])
	tlsConfigMu.Lock()
	defer tlsConfigMu.Unlock()
	if _, exists := registeredTLSConfig[name]; !exists {
		if err := mysql.RegisterTLSConfig(name, &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    roots,
			ServerName: serverName,
		}); err != nil {
			return "", fmt.Errorf("register database TLS configuration: %w", err)
		}
		registeredTLSConfig[name] = struct{}{}
	}
	return name, nil
}
