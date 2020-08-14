package config

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configurations struct {
	Cmd       CmdConfig

}

// Config holds the configuration used for instantiating a new Roach.
type Config struct {
	Database struct {
		// Address that locates our postgres instance
		Host string
		// Port to connect to
		Port string
		// User that has access to the database
		User string
		// Password so that the user can login
		Password string
		// Database to connect to (must have been created priorly)
		DatabaseName string
		// Driver sql driver type
		Driver string
		//SSLModeEnabled dbsslmode disable/enabled
		SSLModeEnabled string
	}
}



type CmdConfig struct {
	AccessToken struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		GrantType        string `json:"grant_type"`
		ClientID         string `json:"client_id"`
		ClientSecret     string `json:"client_secret"`
		Scope            string `json:"scope"`
		SubscriptionID   string `json:"subscription_id"`
		TenantID         string `json:"tenant_id"`
	} `json:"access_token"`
	Resources struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		GrantType        string `json:"grant_type"`
		ClientID         string `json:"client_id"`
		ClientSecret     string `json:"client_secret"`
		Scope            string `json:"scope"`
		SubscriptionID   string `json:"subscription_id"`
		TenantID         string `json:"tenant_id"`
	} `json:"resources"`
	SubscriptionInfo struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		SubscriptionID   string `json:"subscription_id"`
	} `json:"subscriptioninfo"`
	RecommendationList struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"recommendationlist"`
	Recommendation struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"recommendation"`
	ResourceGroups struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"resourcegroups"`
	ResourceGroupCost struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"resourcegroupcost"`
}

// NodeConfig is a struct with the data to be TOML encoded for a node configuration file.
type NodeConfig struct {
	LogToFile       bool   `toml:"log_to_file" json:"log_to_file"`
	LogToStdout     bool   `toml:"log_to_stdout" json:"log_to_stdout"`
	ConfigDirectory string `toml:"config_directory" json:"config_directory"`
	DBDirectory     string `toml:"db_directory" json:"db_directory"`

	WebServer WebServer `toml:"webserver" json:"webserver"`
	Database  Database  `toml:"database" json:"database"`
}

type WebServer struct {
	SiteName        string `toml:"site_name" json:"site_name"`
	SiteDescription string `toml:"site_description" json:"site_description"`
	SiteLatLon      string `toml:"site_lat_lon" json:"site_lat_lon"`
	NodeType        string `toml:"node_type" json:"node_type"`
	ServerIP        string `toml:"server_ip" json:"server_ip"`
	ServerPort      string `toml:"server_port" json:"server_port"`
	GatewayIP       string `toml:"gateway_ip" json:"gateway_ip"`
	GatewayPort     string `toml:"gateway_port" json:"gateway_port"`
	NATSPort        string `toml:"nats_port" json:"nats_port"`
}

type Database struct {
	Host     string `toml:"host" json:"host"`
	Port     string `toml:"port" json:"port"`
	User     string `toml:"user" json:"user"`
	Password string `toml:"password" json:"password"`
	Name     string `toml:"name" json:"name"`
	SslMode  string `toml:"sslmode" json:"sslmode"`
	Driver   string `toml:"driver" json:"driver"`
}

type Server struct {
	ABServerIP               string `toml:"ab_server_ip" json:"server_ip"`
	ABServerPort             string `toml:"ab_server_port" json:"server_port"`
	AzuremonitorFrontEndPort string `toml:"azuremonitor_front_end_port" json:"azuremonitor_front_end_port"`
}

type Debug struct {
	Level int `toml:"debug_level" json:"debug_level"`
}

var debug = Debug{}


func init() {
	debug = LoadDebugConfigs()
}

func GetCmdConfig() (CmdConfig, error) {
	var c CmdConfig
	viper.AddConfigPath("./configs")
	viper.SetConfigName("env.prod")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	filename := viper.ConfigFileUsed()
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, errors.Wrap(err, "read error")
	}
	err = json.Unmarshal([]byte(file), &c)

	if err != nil {
		return c, errors.Wrap(err, "unmarshal")
	}

	return c, nil
}

// Loads configurations for server from proxy.toml
func LoadConfigs() Server {
	ab_server := Server{}

	configPath, err := GetEnvConfigFile()

	if debug.Level > 0 {
		log.Println(configPath)
	}

	if err != nil {
		panic("failed to load .toml config file Configuration file: " + configPath)
	}

	if _, err = toml.DecodeFile(configPath, &ab_server); err != nil {
		panic("failed to load configuration file: " + configPath)
	}

	_, _ = json.Marshal(ab_server)

	return ab_server
}

// Loads Debug type variables from proxy.toml
func LoadDebugConfigs() Debug {
	d := Debug{}

	configPath, err := GetEnvConfigFile()
	if err != nil {
		panic("failed to load .toml config file Configuration file: " + configPath)
	}

	if _, err = toml.DecodeFile(configPath, &d); err != nil {
		panic("failed to load configuration file: " + configPath)
	}

	_, _ = json.Marshal(d)

	return d
}

// LoadNodeConfig loads the TOML file for node configuration settings.
func LoadNodeConfig(file string) NodeConfig {
	nc := NodeConfig{}

	_, err := toml.DecodeFile(file, &nc)

	if err != nil {
		panic("failed to load configuration file: " + file)
	}

	j, _ := json.Marshal(nc)

	log.Println("successfully loaded configuration file: " + file)
	log.Printf(string(j[:]))

	return nc
}

// WriteNodeConfig writes the TOML file for node configuration settings.
func (nc *NodeConfig) WriteNodeConfig(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Fatal("failed to open configuration file")
		return
	}

	defer f.Close()

	err = toml.NewEncoder(f).Encode(nc)

	if err != nil {
		log.Fatal("failed to write configuration to file")
	} else {
		log.Println("configuration written to " + file)
	}
}

// GetDBConfig retrives configuration from toml configs
func GetDBConfig() (Config, error) {
	c := &Config{}

	dbTomlConfig, err := c.LoadDbConfig()
	if err != nil {
		return *c, fmt.Errorf("error: failed to load db config object %v", err)
	}

	c.Database.Driver = dbTomlConfig.Driver
	c.Database.User = dbTomlConfig.User
	c.Database.DatabaseName = dbTomlConfig.Name
	c.Database.Password = dbTomlConfig.Password
	c.Database.Port = dbTomlConfig.Port
	c.Database.SSLModeEnabled = dbTomlConfig.SslMode
	c.Database.Host = dbTomlConfig.Host

	return *c, nil
}

// LoadDbConfig helper function to retrive Database from .toml configurations
// see your directory "env" where all your .toml files are
// located.
func (config *Config) LoadDbConfig() (Database, error) {
	cfx := Database{}

	file, err := GetEnvConfigFile()
	if err != nil {
		panic("failed to load .toml config file Configuration file: " + file)
	}

	//fmt.Printf("Found an Environment file at : %s\n", file)
	_, err = toml.DecodeFile(file, &cfx)

	if err != nil {
		fmt.Printf("error: decoding .toml config file : %+v\n", err)
		panic("failed to Decode .toml Configuration file: " + file)
	}

	_, _ = json.Marshal(cfx)

	return cfx, err
}

// GetEnvConfigFile - returns the .toml configuration based on your enviroment setttings
// projectlevelfolder/env/tomlfile. therefore in dbcontext we going back...back two dirs
func GetEnvConfigFile() (string, error) {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	//return filepath.Join(gopath, "src", "github.com", "aagon00", "azuremonitorServer", "env", "proxy.dev.toml"), nil
	//fmt.Printf("Found GoPath at: %s\n", gopath)

	retVal := ""
	dw, err := os.Getwd()
	fmt.Println(dw)
	if err != nil {
		return retVal, fmt.Errorf("error: os can not read working directory")
	}

	if len(dw) == 0 {
		return retVal, fmt.Errorf("error: directory path can not be empty")
	}

	base := filepath.Base(dw)

	// back...back... (make me smarter please :-)

	//TODO:: proxy.dev.toml Work with team on Configuration Transform patterns #273
	//TODO:: "azuremonitorServer" should not be hardcoded. should be part of a env #273

	switch base {
	case "azuremonitorServer":
		retVal = filepath.Join(dw, "env", "proxy.dev.toml")

	case "dbcontext":
		for i := 0; i < 2; i++ {
			dw = filepath.Dir(dw)
		}
		retVal = filepath.Join(dw, "env", "proxy.dev.toml")

	case "ingest":
		for i := 0; i < 3; i++ {
			dw = filepath.Dir(dw)
		}
		retVal = filepath.Join(dw, "azuremonitor", "env", "proxy.dev.toml")

	case "migration":
		for i := 0; i < 3; i++ {
			dw = filepath.Dir(dw)
		}
		retVal = filepath.Join(dw, "azuremonitor", "env", "proxy.dev.toml")

	//this is reserve for production release
	case "migration scripts":
		fmt.Printf("Deployment migration scripts Job %s\n", dw)
		for i := 0; i < 1; i++ {
			dw = filepath.Dir(dw)
		}
		retVal = filepath.Join(dw, "azuremonitor", "env", "proxy.dev.toml")

	//this is reserve for production release
	case "ingest scripts":

		fmt.Printf("Deployment ingest scripts Job %s\n", dw)
		for i := 0; i < 1; i++ {
			dw = filepath.Dir(dw)
		}
		retVal = filepath.Join(dw, "env", "proxy.dev.toml")

	default:
		return filepath.Join(gopath, "src", "github.com", "Go", "azuremonitor", "env", "proxy.dev.toml"), nil
	}

	//TODO:: Perhaps have 3 levels enum and use that a ref for Information, Errors, Critical etc.
	// for now let's just get some basic going
	if debug.Level > 0 {
		fmt.Printf("The Base directory is: %s\n", base)
		fmt.Printf("Found Working Directory Base: %s\n", dw)
		fmt.Printf("Proposed Config File: %s\n", retVal)
	}

	return retVal, nil
}

// GetConnectionString returns connection string and url connection string
func (config *Config) GetConnectionString() (string, string) {
	// string: connection string.
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password,
		config.Database.DatabaseName, config.Database.SSLModeEnabled)

	// url: postgres://user:password///
	urlConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Database.User,
		config.Database.Password, config.Database.Host, config.Database.Port,
		config.Database.DatabaseName, config.Database.SSLModeEnabled)

	return connectionString, urlConnectionString
}
