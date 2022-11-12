#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE mallbots;
    CREATE USER mallbots_user WITH ENCRYPTED PASSWORD 'mallbots_pass';
    GRANT CONNECT ON DATABASE mallbots TO mallbots_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
    CREATE OR REPLACE FUNCTION created_at_trigger()
    RETURNS TRIGGER AS \$\$
    BEGIN
        NEW.created_at = OLD.created_at;
        RETURN NEW;
    END;
    \$\$ LANGUAGE plpgsql;

    CREATE OR REPLACE FUNCTION updated_at_trigger()
    RETURNS TRIGGER AS \$\$
    BEGIN
        IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
            NEW.updated_at = NOW();
            RETURN NEW;
        ELSE
            RETURN OLD;
        END IF;
    END;
    \$\$ LANGUAGE plpgsql;
EOSQL
