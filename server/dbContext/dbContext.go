package dbContext

import (
	"log"
	"time"
	"os"
	"strconv"
	"strings"
	
	"github.com/Go/server/config"
	"github.com/Go/server/models"
	"github.com/Go/server/util/environment"
	"github.com/Go/server/util/logging"
	"github.com/gocql/gocql"
	hostpool "github.com/hailocab/go-hostpool"
)

var (
	listenAddress string
	webRoot  string
	setupDir string
	Version string
	durationBeforeRetry = time.Second * 20
	retryConnection    = false
	connectionTimeout  = time.Duration(0)
	keepalive          = time.Duration(0)
	reconnect          = time.Duration(0)
	clusterConsistency = gocql.LocalQuorum
	clusterDC          = os.Getenv(config.CassandraDCEnv)
	queryRetries       int
	cql  *gocql.Session
)

func InitDBConnection() {
	var err error
	cql, err = getCassandraSession(config.RecordKeyspace)
	if err != nil {
		logging.Fatalf("Unable to establish DB Connection %s, Error: %+v.", config.RecordKeyspace, err)
	}
}

func getCassandraSession(keyspace string) (*gocql.Session, error) {
	var cassandraHosts = make([]string, 0)
	cassandraHostList := environment.GetSet(config.CassandraNodesEnv, "127.0.0.1")
	for _, host := range strings.Split(cassandraHostList, ",") {
		cassandraHosts = append(cassandraHosts, host)
	}
	return OpenClusterConnection(keyspace, cassandraHosts)
}
// OpenClusterConnection opens a cluster connection for the specified keySpace.
func OpenClusterConnection(keyspace string, clusterNodes []string) (*gocql.Session, error) {
	numConnsString := environment.GetSet(config.CassandraNumConnsEnv, "1")
	numConns, err := strconv.Atoi(numConnsString)
	if err != nil {
		logging.Fatalf("Unable to parse %q env variable with entry: %q: %+v", config.CassandraNumConnsEnv, numConnsString, err)
	}
	return OpenClusterConnectionConns(keyspace, numConns, clusterNodes)
}

//OpenClusterConnectionConns opens a session with the specified number of connections.
//Most applications should just ust OpenClusterConnection
func OpenClusterConnectionConns(keyspace string, numConns int, clusterNodes []string) (*gocql.Session, error) {
	handleEnvVars()
	clusterCfg := getClusterCfg(keyspace, numConns, clusterNodes)
	return createSession(clusterCfg)
}

// Spin forever until connection to cassandra is established, unless the
// environment variable RETRY_CASSANDRA_CONNECTION is set to 0.
func createSession(clusterCfg *gocql.ClusterConfig) (*gocql.Session, error) {
	for {
		session, err := clusterCfg.CreateSession()
		if err != nil {
			if !retryConnection {
				log.Fatalf("Error creating cluster session: %v", err)
			}
			log.Println("Cluster connection keyspace '", clusterCfg.Keyspace,
				"' clusterNodes '", clusterCfg.Hosts, "' err =", err)
			time.Sleep(durationBeforeRetry)
		} else {
			return session, err
		}
	}
}

// Retrieve environment variables if available.
func handleEnvVars() {
	retryConnectionStr := os.Getenv(config.RetryCassandraConnection)
	if retryConnectionStr == "1" || retryConnectionStr == "true" {
		retryConnection = true
	}

	timeoutSecStr := os.Getenv(config.CassandraConnectionTimeoutSec)
	if len(timeoutSecStr) == 0 {
		timeoutSecStr = "30.0"
	}
	timeoutSec, err := strconv.ParseFloat(timeoutSecStr, 64)
	if err == nil && timeoutSec > 0.0 {
		connectionTimeout = time.Millisecond *
			time.Duration(timeoutSec*config.SecondToMillisecond)
	} else if err != nil {
		log.Fatalf("%s is in wrong format.  Expecting a float value", config.CassandraConnectionTimeoutSec)
	}
	keepalive, err = time.ParseDuration(environment.GetSet(config.CassandraKeepaliveEnv, "0s"))
	if err != nil {
		log.Fatalf("%s is in the wrong format.  Expecting a number in the form 0s", config.CassandraKeepaliveEnv)
	}

	reconnect, err = time.ParseDuration(environment.GetSet(config.CassandraReconnectEnv, "0s"))
	if err != nil {
		log.Fatalf("%s is in the wrong format.  Expecting a number in the form 0s", config.CassandraReconnectEnv)
	}

	consistency := os.Getenv(config.CassandraConsistency)
	consistency = strings.ToUpper(consistency)
	switch consistency {
	case "ANY":
		clusterConsistency = gocql.Any
	case "ONE":
		clusterConsistency = gocql.One
	case "TWO":
		clusterConsistency = gocql.Two
	case "THREE":
		clusterConsistency = gocql.Three
	case "QUORUM":
		clusterConsistency = gocql.Quorum
	case "ALL":
		clusterConsistency = gocql.All
	case "LOCALQUORUM":
		clusterConsistency = gocql.LocalQuorum
	case "EACHQUORUM":
		clusterConsistency = gocql.EachQuorum
	case "LOCALONE":
		clusterConsistency = gocql.LocalOne
	default:
		clusterConsistency = gocql.LocalQuorum
	}

	queryRetriesStr := os.Getenv(config.CassandraQueryRetriesEnv)
	if queryRetriesStr == "" {
		queryRetriesStr = "0"
	}
	queryRetries, err = strconv.Atoi(queryRetriesStr)
	if err != nil {
		log.Fatalf("%s is in wrong format.  Expecting an int value", config.CassandraQueryRetriesEnv)
	}
}

func getClusterCfg(keyspace string, numConns int, clusterNodes []string) *gocql.ClusterConfig {
	clusterCfg := gocql.NewCluster(clusterNodes...)
	clusterCfg.ProtoVersion = 3
	clusterCfg.Consistency = clusterConsistency
	clusterCfg.Compressor = gocql.SnappyCompressor{}
	if connectionTimeout > time.Duration(0) {
		clusterCfg.Timeout = connectionTimeout
		clusterCfg.ConnectTimeout = connectionTimeout
	}
	clusterCfg.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: queryRetries}
	clusterCfg.SocketKeepalive = keepalive
	clusterCfg.ReconnectInterval = reconnect
	if numConns < 1 {
		numConns = 1
	}
	clusterCfg.NumConns = numConns
	clusterCfg.Keyspace = keyspace
	if clusterDC != "" {
		clusterCfg.HostFilter = gocql.DataCentreHostFilter(clusterDC)
	}
	clusterCfg.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(hostpool.New(nil))
	return clusterCfg
}

func GetRecordsByHash(hashValue string) ([]models.Record, error) {
	var records = make([]models.Record, 0)
	r := models.Record{}
	srcGroups := make([]string, 0)
	trgGroups := make([]string, 0)
	var iter = cql.Query(config.GetAllRecords).Iter()
	for iter.Scan(&r.Source.ID, &r.LinkType, &r.Target.ID, &r.Source.Data, &srcGroups, &trgGroups, &r.Target.Data ) {
		if r.LinkType == "files" || r.LinkType == "selectors" {
			if strings.Contains(r.Target.Data,hashValue) {
				records = append(records, r)
				r = models.Record{}
			}
		}
	}
	return records, iter.Close()
}

func GetTargetData(sourceID string) ([]models.Record, error) {
	var records = make([]models.Record, 0)
	r := models.Record{}
	groups := make([]string, 0)
	var iter = cql.Query(config.GetTargetDataCQL, sourceID).Iter()
	for iter.Scan(&groups, &r.Target.Data) {
		records = append(records, r)
		r = models.Record{}
	}
	return records, iter.Close()
}

//NewTargetData creates a new record if it doesn't exist,
//otherwise it returns the old record.
func NewTargetData(r models.Record) (models.Record, error) {
	mapCAS := map[string]interface{}{}
	applied, err := cql.Query(config.InsertDataIfNotExistsCQL,
		r.Source.ID,
		r.Source.Groups,
		r.Source.Data,
		r.LinkType,
		r.Target.ID,
		r.Target.Groups,
		r.Target.Data).MapScanCAS(mapCAS)
		
	if applied {
		return r, err
	} else {
		var curVal models.Record
		curVal.Target.ID, _ = mapCAS["targetid"].(string)
		curVal.Target.Data, _ = mapCAS["targetdata"].(string)
		strAdg, _ := mapCAS["targetadgroups"].([]string)
		curVal.Target.DataSecurity.Groups = strAdg
		return curVal, err
	}
}

//GetTargetDataByLinkType returns target data by linktype.
func GetTargetDataByLinkType(sourceID, linkType string) ([]models.Record, error) {
	var records = make([]models.Record, 0)
	r := models.Record{}
	groups := make([]string, 0)
	var iter = cql.Query(config.GetTargetDataByLinkTypeCQL, sourceID, linkType).Iter()
	for iter.Scan(&groups, &r.Target.ID, &r.Target.Data) {
		records = append(records, r)
		r = models.Record{}
	}

	return records, iter.Close()
}

//GetStaticCols returns the static columns of a sourceid
func GetStaticCols(sourceID string) (models.Data, error) {
	var src models.Data
	var ret models.Data
	var err error
	src.ID = sourceID
	err = nil
	iter := cql.Query(config.GetStaticColsCQL, sourceID).Iter()
	for iter.Scan(&src.Classification, &src.Groups, &src.Data) {
		ret = src
	}
	if err == nil {
		err = iter.Close()
	} else {
		iter.Close()
	}
	return ret, err
}

//GetRecordsBySource returns all records that have a specific sourceid
func GetRecordsBySource(sourceID string) ([]models.Record, error) {
	var records = make([]models.Record, 0)
	r := models.Record{}
	srcGroups := make([]string, 0)
	trgGroups := make([]string, 0)
	var iter = cql.Query(config.GetRecordsBySourceID, sourceID).Iter()
	for iter.Scan(&r.Source.ID, &r.LinkType, &r.Target.ID, &r.Source.Data, &srcGroups, &trgGroups, &r.Target.Data) {
		records = append(records, r)
		r = models.Record{}
	}
	return records, iter.Close()
}

//UpdateTargetData updates a record.
func UpdateTargetData(r models.Record) (models.Record, error) {
	err := cql.Query(config.InsertDataCQL,
		r.Source.ID,
		r.Source.Groups,
		r.Source.Data,
		r.LinkType,
		r.Target.ID,
		r.Target.Groups,
		r.Target.Data).Exec()
	if err != nil {
		return r, err
	}
	return r, nil
}

//DeleteData deletes the record.
func DeleteData(sourceID, linkType, targetID string) error {
	return cql.Query(config.DeleteDataCQL, sourceID, linkType, targetID).Exec()
}


//UpdateRecordsByID finds all records where sourceID=data.ID and updates the
// data there, then it switches it around and updates the reverse link
// linktype should be as if the data.ID where the target
// for example, if data.ID is a batchID then linktype=BatchLinkType
func UpdateRecordsByID(data models.Data, linkType string) error {
	srcRecords, err := GetRecordsBySource(data.ID)
	if err != nil {
		return err
	}
	for _, cur := range srcRecords {
		cur.Source = data
		_, err = UpdateTargetData(cur)
		if err != nil {
			return err
		}
		if cur.LinkType != "" {
			rev := models.Record{
				Source:   cur.Target,
				LinkType: linkType,
				Target:   cur.Source,
			}
			_, err = UpdateTargetData(rev)

			if err != nil {
				return err
			}
		}
	}
	return err
}
