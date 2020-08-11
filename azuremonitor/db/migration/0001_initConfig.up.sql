SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" schema pg_catalog version "1.1";

DROP SCHEMA IF EXISTS azuremonitor CASCADE; 
CREATE SCHEMA IF NOT EXISTS azuremonitor;
ALTER SCHEMA azuremonitor OWNER TO postgres;  

-- added to support reset migration script
CREATE TABLE IF NOT EXISTS azuremonitor.schema_migrations
(
    version bigint NOT NULL,
    CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
)
    
TABLESPACE pg_default;
ALTER TABLE azuremonitor.schema_migrations
OWNER to postgres;

SET search_path = azuremonitor, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;
