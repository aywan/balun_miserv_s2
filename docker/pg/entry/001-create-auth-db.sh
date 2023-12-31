#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE ${AUTH_DB_NAME};
	CREATE ROLE ${AUTH_DB_USER} WITH LOGIN PASSWORD '${AUTH_DB_PASS}';
	GRANT ALL ON DATABASE ${AUTH_DB_NAME} TO ${AUTH_DB_USER};
	ALTER DATABASE ${AUTH_DB_NAME} OWNER TO ${AUTH_DB_USER};
EOSQL
