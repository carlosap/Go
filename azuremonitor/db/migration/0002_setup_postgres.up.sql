DO
$body$
BEGIN
-- user to login and 
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_user WHERE usename = 'azmonitor'
    ) THEN
        CREATE ROLE azmonitor WITH LOGIN;
        GRANT USAGE ON SCHEMA azmonitor TO azmonitor;
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_group WHERE groname = 'azmonitorfuncs'
    ) THEN
        CREATE ROLE azmonitorfuncs WITH NOLOGIN;
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_group WHERE groname  = 'azmonitorread'
    ) THEN
        CREATE ROLE azmonitorread WITH NOLOGIN;
    END IF;  
END
$body$
;
