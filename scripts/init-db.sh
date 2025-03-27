#!/bin/bash
set -e

# Wait for PostgreSQL to start
echo "Waiting for PostgreSQL to start..."
sleep 5

# Create the database and user if they don't exist
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE $POSTGRES_DB;
    CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';
    GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;
EOSQL

echo "Database initialized successfully!"