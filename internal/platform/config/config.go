// Package config validates process configuration without logging secret values.
package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Address, DatabaseURL, Environment string
	Docs                              bool
	Level                             slog.Level
	DBTimeout, ShutdownTimeout        time.Duration
}

func Load(get func(string) string) (Config, error) {
	c := Config{Address: "127.0.0.1:8080", Environment: "local", DBTimeout: 2 * time.Second, ShutdownTimeout: 10 * time.Second}
	if v := get("APP_ENV"); v != "" {
		c.Environment = v
	}
	if c.Environment != "local" && c.Environment != "test" && c.Environment != "production" {
		return c, fmt.Errorf("APP_ENV must be local, test, or production")
	}
	if v := get("HTTP_ADDR"); v != "" {
		c.Address = v
	}
	host, port, err := net.SplitHostPort(c.Address)
	if err != nil {
		return c, fmt.Errorf("HTTP_ADDR must be host:port")
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return c, fmt.Errorf("HTTP_ADDR port must be 1..65535")
	}
	if c.Environment == "local" && host != "localhost" && (net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback()) {
		return c, fmt.Errorf("local HTTP_ADDR must be loopback")
	}
	c.DatabaseURL = get("DATABASE_URL")
	u, err := url.Parse(c.DatabaseURL)
	if err != nil || u == nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Hostname() == "" || u.User == nil || u.User.Username() == "" || u.Path == "" || u.Path == "/" {
		return c, fmt.Errorf("DATABASE_URL must be a PostgreSQL URL with user, host, and database")
	}
	if c.Environment == "production" && u.Query().Get("sslmode") != "verify-full" {
		return c, fmt.Errorf("production DATABASE_URL requires sslmode=verify-full")
	}
	if v := get("DOCS_ENABLED"); v != "" {
		c.Docs, err = strconv.ParseBool(v)
		if err != nil {
			return c, fmt.Errorf("DOCS_ENABLED must be boolean")
		}
	}
	if c.Docs && c.Environment != "local" {
		return c, fmt.Errorf("docs are restricted to local environment")
	}
	if v := get("LOG_LEVEL"); v != "" {
		if err := c.Level.UnmarshalText([]byte(v)); err != nil {
			return c, fmt.Errorf("LOG_LEVEL is invalid")
		}
	}
	for key, target := range map[string]*time.Duration{"DB_TIMEOUT": &c.DBTimeout, "SHUTDOWN_TIMEOUT": &c.ShutdownTimeout} {
		if v := get(key); v != "" {
			d, e := time.ParseDuration(v)
			if e != nil || d <= 0 || d > time.Minute {
				return c, fmt.Errorf("%s must be positive and at most 1m", key)
			}
			*target = d
		}
	}
	return c, nil
}
