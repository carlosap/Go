-- // create Table
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS azmonitor.application (
  applicationID SERIAL,
  subscription_id varchar,
  name varchar,
  tenant_id varchar,
  grant_type varchar,
  client_id varchar,
  client_secret varchar,
  lastmodified TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (applicationID)
);
