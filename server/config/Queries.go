package config

const (	
	GetTargetDataByLinkTypeCQL = "SELECT targetadgroups, targetid, targetdata FROM records WHERE sourceid = ? AND linktype =?"
	GetStaticColsCQL           = "SELECT distinct sourceclassification, sourceadgroups, sourcedata from records where sourceid = ?"
	GetRecordsBySourceID       = "SELECT sourceid, linktype, targetid, sourcedata, sourceadgroups, targetadgroups, targetdata FROM records WHERE sourceid = ?"
	GetAllRecords              = "SELECT sourceid, linktype, targetid, sourcedata, sourceadgroups, targetadgroups, targetdata FROM records"
	GetTargetDataCQL           = "SELECT targetadgroups, targetdata FROM records WHERE sourceid = ?"
	GetTargetSample            = "SELECT targetadgroups, targetdata FROM records LIMIT 10"
	InsertDataCQL              = "INSERT INTO records (sourceid, sourceadgroups, sourcedata, linktype, targetid, targetadgroups, targetdata) VALUES(?,?,?,?,?,?,?)"
	InsertDataIfNotExistsCQL   = "INSERT INTO records (sourceid, sourceadgroups, sourcedata, linktype, targetid, targetadgroups, targetdata) VALUES(?,?,?,?,?,?,?) IF NOT EXISTS"
	DeleteDataCQL              = "DELETE FROM records WHERE sourceid = ? AND linktype = ? AND targetid = ?"
)