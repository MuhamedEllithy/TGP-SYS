// dbcheck verifies the exact application's host TCP path without revealing credentials.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gateplaza.local/facility/internal/platform/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func main() {
	test := flag.Bool("test", false, "verify only the dedicated disposable Phase 0 database")
	absent := flag.Bool("probe-absent", false, "assert probe rollback on disposable DB")
	flag.Parse()
	if *absent && !*test {
		fmt.Fprintln(os.Stderr, "probe rollback check requires --test")
		os.Exit(1)
	}
	if err := check(*test, *absent); err != nil {
		// Driver errors/URLs must never be printed: they can contain passwords.
		fmt.Fprintln(os.Stderr, "Database verification failed; check local settings, server authentication and project isolation without displaying credentials.")
		os.Exit(1)
	}
}

func check(test, probeAbsent bool) error {
	raw, expectedDB, expectedUser := os.Getenv("DATABASE_URL"), os.Getenv("DB_NAME"), os.Getenv("DB_USER")
	if test {
		raw, expectedDB, expectedUser = os.Getenv("TEST_DATABASE_URL"), "tgp_facility_phase0_test", os.Getenv("MIGRATION_USER")
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil || parsed.Hostname() != "127.0.0.1" || parsed.Port() != os.Getenv("DB_PORT") || expectedDB == "" || expectedUser == "" {
		return errors.New("invalid dedicated local connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := db.Open(ctx, raw)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return err
	}
	var name, user, version string
	var super, createdb, createrole, createSchema bool
	err = pool.QueryRow(ctx, `SELECT current_database(), current_user, current_setting('server_version'),
  rolsuper, rolcreatedb, rolcreaterole, has_schema_privilege(current_user, 'public', 'CREATE')
  FROM pg_roles WHERE rolname = current_user`).Scan(&name, &user, &version, &super, &createdb, &createrole, &createSchema)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(version, "17.") || name != expectedDB || user != expectedUser || super || createdb || createrole || (!test && createSchema) {
		return errors.New("identity or privilege mismatch")
	}
	if test {
		var unexpected int
		err = pool.QueryRow(ctx, `SELECT count(*) FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace
   WHERE n.nspname NOT LIKE 'pg_%' AND n.nspname <> 'information_schema'
   AND c.relkind IN ('r','p','v','m','S','f')
   AND NOT (n.nspname = 'public' AND c.relname IN ('platform_probe','schema_migrations'))`).Scan(&unexpected)
		if err != nil || unexpected != 0 {
			return errors.New("disposable database contains unexpected objects")
		}
	}
	if probeAbsent {
		var absent bool
		if err := pool.QueryRow(ctx, "SELECT to_regclass('public.platform_probe') IS NULL").Scan(&absent); err != nil || !absent {
			return errors.New("probe rollback not established")
		}
	}
	// A successful connection followed by an authentication-specific rejection with
	// a fresh wrong password proves this host path isn't silently accepting trust auth.
	cfg, err := pgx.ParseConfig(raw)
	if err != nil {
		return err
	}
	token := make([]byte, 32)
	if _, err = rand.Read(token); err != nil {
		return err
	}
	cfg.Password = hex.EncodeToString(token)
	bad, err := pgx.ConnectConfig(ctx, cfg)
	if err == nil {
		_ = bad.Close(ctx)
		return errors.New("password authentication was not enforced")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "28P01" {
		return errors.New("wrong-password authentication rejection not established")
	}
	fmt.Printf("PASS: authenticated host TCP connection; database=%s user=%s PostgreSQL=%s; non-superuser; wrong password rejected\n", name, user, version)
	return nil
}
