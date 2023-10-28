#!/usr/bin/env bash

set -x

MIGRATION_DSN="host=${DB_HOST} port=${DB_PORT:-5432} dbname=${DATABASE} user=${DB_USER} password='${DB_PASS}' sslmode=${DB_SSLMODE:-disable}"
sleep 2 && goose postgres "${MIGRATION_DSN}" up -v
