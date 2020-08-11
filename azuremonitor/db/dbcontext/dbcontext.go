/*Todo: notes
//https://www.postgresql.org/docs/9.4/functions-json.html
//https://godoc.org/github.com/go-pg/pg
*/

package dbcontext

import (
	"database/sql"

	"github.com/Go/azuremonitor/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// DbContext Prefer running queries from DB unless there is a specific need for a continuous single database connection
// DB - is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines.
// Conn - represents a single database connection rather than a pool of database connections.
type DbContext struct {
	Pgdb   *pg.DB
	Sqldb  *sql.DB
	Config config.Config
}

// NewDbConctext starts a new hanlder/instance
func NewDbContext() (*DbContext, error) {
	dbcontext := &DbContext{}

	cfg, _ := config.GetDBConfig()
	if cfg.Database.Host == "" || cfg.Database.Port == "" || cfg.Database.User == "" ||
		cfg.Database.Password == "" || cfg.Database.DatabaseName == "" {
		err := errors.Errorf(
			"All fields must be set (%s)",
			spew.Sdump(cfg))
		return dbcontext, err
	}

	_, connectionString := cfg.GetConnectionString()

	sqldb, err := sql.Open(cfg.Database.Driver, connectionString)

	if err != nil {
		return nil, err
	}

	options, err := pg.ParseURL(connectionString)

	if err != nil {
		return nil, err
	}
	pgdb := pg.Connect(options)

	dbcontext.Pgdb = pgdb
	dbcontext.Sqldb = sqldb
	dbcontext.Config = cfg

	return dbcontext, nil
}

// Close ensures both connections are close
func (d *DbContext) Close() (err error) {

	_ = d.Pgdb.Close()
	_ = d.Sqldb.Close()
	return
}

//============================================Common Driver Requirements===================================
//Columns matrix for querying columns db
var Columns = struct {
	Application struct {
		Applicationid, SubscriptionID, Name, TenantID, GrantType, ClientID, ClientSecret, Lastmodified string
	}
}{
	Application: struct {
		Applicationid, SubscriptionID, Name, TenantID, GrantType, ClientID, ClientSecret, Lastmodified string
	}{
		Applicationid:  "applicationid",
		SubscriptionID: "subscription_id",
		Name:           "name",
		TenantID:       "tenant_id",
		GrantType:      "grant_type",
		ClientID:       "client_id",
		ClientSecret:   "client_secret",
		Lastmodified:   "lastmodified",
	},
}

// Tables matrix for querying tables db
var Tables = struct {
	Application struct {
		Name string
	}
}{
	Application: struct {
		Name string
	}{
		Name: "azmonitor.application",
	},
}
