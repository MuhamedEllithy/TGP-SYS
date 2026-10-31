#!/bin/sh
set -eu
# URL comes from container environment, never a host CLI argument.
exec migrate -path /migrations -database "$MIGRATION_DATABASE_URL" "$@"
