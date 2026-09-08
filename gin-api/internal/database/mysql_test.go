package database

import (
	"strings"
	"testing"
)

func TestMySQLDSNConvertsAivenURL(t *testing.T) {
	dsn, err := mysqlDSN("mysql://avnadmin:p%40ss@example.aivencloud.com:11881/defaultdb?ssl-mode=REQUIRED")
	if err != nil {
		t.Fatalf("mysqlDSN returned an error: %v", err)
	}
	for _, expected := range []string{
		"avnadmin:p@ss@tcp(example.aivencloud.com:11881)/defaultdb",
		"tls=true",
		"parseTime=true",
	} {
		if !strings.Contains(dsn, expected) {
			t.Fatalf("expected DSN to contain %q, got %q", expected, dsn)
		}
	}
}

func TestMySQLDSNRejectsConnectionWithoutSSLMode(t *testing.T) {
	_, err := mysqlDSN("mysql://user:password@example.com:3306/database")
	if err == nil {
		t.Fatal("expected connection without ssl-mode to be rejected")
	}
}

func TestMySQLDSNAcceptsDisabledTLSForLocalDevelopment(t *testing.T) {
	dsn, err := mysqlDSN("mysql://app:password@db:3306/digital_stamp?ssl-mode=DISABLED")
	if err != nil {
		t.Fatalf("mysqlDSN returned an error: %v", err)
	}
	if strings.Contains(dsn, "tls=true") {
		t.Fatalf("expected local DSN not to enable TLS, got %q", dsn)
	}
}
