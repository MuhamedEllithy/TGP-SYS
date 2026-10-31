//go:build integration

package integration

import (
	"context"
	"github.com/jackc/pgx/v5"
	"net/url"
	"os"
	"testing"
	"time"
)

// Run only against a disposable database after migrate-up. This suite does not
// reset schemas or execute destructive DDL on an arbitrary DATABASE_URL.
func TestDatabaseProbe(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	u, e := url.Parse(raw)
	if e != nil || u == nil || (u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost") {
		t.Fatal("TEST_DATABASE_URL must explicitly target local disposable PostgreSQL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, e := pgx.Connect(ctx, raw)
	if e != nil {
		t.Fatal("database connection failed")
	}
	defer conn.Close(ctx)
	if e := conn.Ping(ctx); e != nil {
		t.Fatal("database ping failed")
	}
	var name string
	if e := conn.QueryRow(ctx, "SELECT current_database()").Scan(&name); e != nil || name != "tgp_facility_phase0_test" {
		t.Fatal("refusing a database other than the disposable Phase 0 target")
	}
	var count int
	if e := conn.QueryRow(ctx, "SELECT count(*) FROM platform_probe WHERE singleton").Scan(&count); e != nil || count != 1 {
		t.Fatal("migration probe missing")
	}
	tx, e := conn.Begin(ctx)
	if e != nil {
		t.Fatal("transaction initialization failed")
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "INSERT INTO platform_probe(singleton) VALUES(false)"); e == nil {
		t.Fatal("probe CHECK constraint not enforced")
	}
}
