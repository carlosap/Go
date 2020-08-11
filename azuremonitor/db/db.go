package db

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/Go/azuremonitor/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const sqlDriverName string = "postgres"

// Instance holds a connect pool to the database
type Database struct {
	Log    *logrus.Logger
	Config *config.Database
	Db     *sql.DB
}

// NewDatabase will create an initial instance for the database
func NewDatabase(log *logrus.Logger, config *config.Database) *Database {
	db, err := OpenConnection(config)

	if err != nil {
		log.WithError(err).Panic("failed to open database")
	}

	// Ping db to test connection
	if err = db.Ping(); err != nil {
		log.WithError(err).Panic("failed to ping database")
	}

	return &Database{Log: log, Config: config, Db: db}
}

// Close all resources used used upon exit or panic
func (db *Database) Close() {
	if db.Db == nil {
		return
	}

	if err := db.Db.Close(); err != nil {
		db.Log.WithError(err).Error("failed to close database")
	}
}

func OpenConnection(config *config.Database) (*sql.DB, error) {
	return sql.Open(sqlDriverName, BuildConnectionInfo(config))
}

func BuildConnectionInfo(config *config.Database) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.User,
		config.Password,
		config.Name,
		config.Host,
		config.Port)
}

func checkPrependComma(buffer *bytes.Buffer) {
	if buffer.Len() > 0 {
		buffer.WriteString(", ")
	}
}

func writeColumn(buffer *bytes.Buffer, String string) {
	checkPrependComma(buffer)
	buffer.WriteString(String)
}

func writeValue(buffer *bytes.Buffer, String string) {
	checkPrependComma(buffer)
	buffer.WriteString(fmt.Sprintf("'%s'", String))
}
