package config

const (
	RecordKeyspace = "texas"
	RecordKeyspaceEnv = "RECORD_KEYSPACE"
	CassandraNodesEnv = "CASSANDRA_NODES"
	CassandraDCEnv = "CASSANDRA_DATACENTER"
	RetryCassandraConnection = "RETRY_CASSANDRA_CONNECTION"
	CassandraNumConnsEnv = "CASSANDRA_NUM_CONNS"
	CassandraConsistency = "CASSANDRA_CONSISTENCY"	
	CassandraKeepaliveEnv = "CASSANDRA_KEEPALIVE"
	CassandraReconnectEnv = "CASSANDRA_RECONNECT"
	CassandraQueryRetriesEnv = "CASSANDRA_QUERY_RETRIES"
	CassandraConnectionTimeoutSec = "CASSANDRA_CONNECTION_TIMEOUT_SEC"
	SecondToMillisecond = 1000
)