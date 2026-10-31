#!/bin/sh
# Official-image first-boot hook. Never runs against a reused initialized volume.
set -eu
export PGPASSWORD="$POSTGRES_PASSWORD"
# Do not enable shell tracing, echo variables, or log SQL containing passwords.
if ! psql --no-psqlrc --username "$POSTGRES_USER" --dbname postgres --quiet --set ON_ERROR_STOP=1 > /dev/null 2>&1 <<'SQL'
SET log_statement = 'none';
SET log_min_error_statement = 'panic';
\getenv app_user DB_USER
\getenv app_password DB_PASSWORD
\getenv migration_user MIGRATION_USER
\getenv migration_password MIGRATION_PASSWORD
\getenv dev_db DB_NAME
\getenv test_db TEST_DB_NAME
CREATE ROLE :"migration_user" LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD :'migration_password';
CREATE ROLE :"app_user" LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD :'app_password';
CREATE DATABASE :"dev_db" OWNER :"migration_user";
CREATE DATABASE :"test_db" OWNER :"migration_user";
REVOKE ALL ON DATABASE :"dev_db" FROM PUBLIC;
REVOKE ALL ON DATABASE :"test_db" FROM PUBLIC;
GRANT CONNECT ON DATABASE :"dev_db" TO :"app_user";
GRANT CONNECT ON DATABASE :"test_db" TO :"app_user";
\connect :dev_db
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO :"app_user";
GRANT USAGE, CREATE ON SCHEMA public TO :"migration_user";
\connect :test_db
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO :"app_user";
GRANT USAGE, CREATE ON SCHEMA public TO :"migration_user";
SQL
then
    echo 'TGP-Sys first-boot database initialization failed; investigate without displaying credentials.' >&2
    exit 1
fi
unset PGPASSWORD
