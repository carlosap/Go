-- // create Table
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS azuremonitor.azuremonitor (
  azuremonitor_id UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
  name varchar,
  hostname int,
  lastmodified TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (azuremonitor_id)
);
