DO
$body$
BEGIN
-- user to login and 
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_user WHERE usename = 'azuremonitor'
    ) THEN
        CREATE ROLE azuremonitor WITH LOGIN;
        GRANT USAGE ON SCHEMA azuremonitor TO azuremonitor;
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_group WHERE groname = 'azuremonitorfuncs'
    ) THEN
        CREATE ROLE azuremonitorfuncs WITH NOLOGIN;
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_group WHERE groname  = 'azuremonitorread'
    ) THEN
        CREATE ROLE azuremonitorread WITH NOLOGIN;
    END IF;  
END
$body$
;
