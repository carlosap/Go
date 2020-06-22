package main

import (
	"flag"
	//"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/aagon00/POILibs/file"
	"github.com/aagon00/POILibs/networking"
	"github.com/sirupsen/logrus"
)

//var log = logrus.New()
var instanceConfig = EphIConfig{}
var tunnelOptions = TunnelOptions{}
var regionOptions = RegionZonesOptions{}

var Selections PoPs

type PoPInfo struct {
	Name             string
	ConfFilename     string
	SecretsFileName  string
	PoPCloudProvider string
	PoPIP            string
	EntryIP          string
	Username         string
}

type PoPs struct {
	Option []PoPInfo
}

/*
 - this should be removed; it was a quick fix for "func networkInfo" in handler.go
*/
func PollStatus() bool {
	return true
}

/*
	function that stores the username in the tunnel config
	NOTE: review the configs file to see the sturctures being used
*/
func setUsername(_un string) {
	/* can not leave 0 as index */
	tunnelOptions.TunnelOptions[0].Username = _un
}

/* array of objects that contain all region information */
func getRegionsAndZones() []RegionZone {
	return regionOptions.RegionsAndZones
}

func getDirectory() string { // gets directory info
	dir, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}
	dir = strings.Replace(dir, "server", "", 1)
	dir = dir + "hf_UI/build/"
	return dir
}

func CheckInterfaceStatus(ifName string) bool { // Checks interface status
	sysIf, err1 := net.InterfaceByName(ifName)
	if err1 == nil {
		addrList, err2 := sysIf.Addrs()
		if err2 == nil {
			if len(addrList) > 0 {
				return true
			}
			if len(addrList) == 0 {
				return false
			}
		}
	}
	return false
}

/*
	tunnel params contains only  POP, POE and Name

NOTE:
	this function should start to generate the tunnel
	should return a boolean?
*/
func setCustomTunnelConfig(username string, hops int, tunnelParams []string) int {
	rand.Seed(693)

	var entryOption RegionZone
	var popOption RegionZone
	for _, tunOption := range regionOptions.RegionsAndZones {
		if hops > 1 {
			if tunOption.RegionName == tunnelParams[0] {
				entryOption = tunOption
			}
		}
		if tunOption.RegionName == tunnelParams[1] {
			popOption = tunOption
		}
	}

	popIndex := 0
	if hops > 1 {
		popIndex = 1
	}
	customConfig := TunnelConfig{Name: tunnelParams[2], Username: username, PointOfEntry: 0, PointOfPresence: popIndex}
	// TODO: check the cloud provider in the inputs to sanitize for valid provider types...
	if hops > 1 {
		var entryCloudConfig interface{}
		entryZone := ""
		entryImageName := ""
		if entryOption.Provider == "aws" {
			entryCloudConfig = AWSConfig{}
			entryImageName = "arch linux"
		}
		if entryOption.Provider == "azure" {
			entryCloudConfig = AzureConfig{}
			entryImageName = "debian linux"
		}
		if entryOption.Provider == "gce" {
			entryCloudConfig = GCEConfig{}
			entryImageName = "debian linux"
		}
		if len(entryOption.Zones) == 1 {
			entryZone = entryOption.Zones[0]
		} else {
			entryZone = entryOption.Zones[rand.Intn(len(entryOption.Zones))]
		}
		entryNode := NodeConfig{CloudProvider: entryOption.Provider, CloudConfig: entryCloudConfig, Region: entryZone, ImageName: entryImageName, PrivateSubnet: networking.GenerateIP() + "/24", VirtualSubnet: networking.GenerateIP() + "/28", SSHPort: networking.GeneratePort()}
		customConfig.Nodes = append(customConfig.Nodes, entryNode)
	}
	var popCloudConfig interface{}
	popZone := ""
	popImageName := ""
	if popOption.Provider == "aws" {
		popCloudConfig = AWSConfig{}
		popImageName = "arch linux"
	}
	if popOption.Provider == "azure" {
		popCloudConfig = AzureConfig{}
		popImageName = "debian linux"
	}
	if popOption.Provider == "gce" {
		popCloudConfig = GCEConfig{}
		popImageName = "debian linux"
	}
	if len(popOption.Zones) == 1 {
		popZone = popOption.Zones[0]
	} else {
		popZone = popOption.Zones[rand.Intn(len(popOption.Zones))]
	}
	virtSubnet := ""
	if hops == 1 {
		virtSubnet = networking.GenerateIP() + "/28"
	}
	popNode := NodeConfig{CloudProvider: popOption.Provider, CloudConfig: popCloudConfig, Region: popZone, ImageName: popImageName, PrivateSubnet: networking.GenerateIP() + "/24", VirtualSubnet: virtSubnet, SSHPort: networking.GeneratePort(), AllowAllTrafficOut: true}
	customConfig.Nodes = append(customConfig.Nodes, popNode)

	// TODO: decide whether to force user to decide on regions for nodes above 2 or random select...
	// TODO: fix random port generator in networking package to be in certain range...
	// TODO: decide whether to return the struct or call cloud and pass from here...
	// TODO: begin work here calling cloud loader for selected PoP...
	return 1 //success notify to browser
}

/*******************************/
/*
	below can be removed once server is in main HaroldFinch repo with CLI tool
*/

// initDebugLog instantiates the debug logging system, based on flags for file and stdout will log output accordingly.
func initDebugLog() {
	//var multiLogger io.Writer

	//logFile, err := os.OpenFile(filepath.Join(instanceConfig.LogDirectory+instanceConfig.Name+"_log_debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{"Function": "initDebugLog", "File": instanceConfig.Name + "_log_debug.log", "Message": "Failed to open debug log file for server instance."}).Fatal(err)
	//}

	//multiLogger = io.MultiWriter(logFile, os.Stdout)

	//if !instanceConfig.LogToFile && instanceConfig.LogToStdout {
	//	logrus.Out = os.Stdout
	//} else if instanceConfig.LogToFile && !instanceConfig.LogToStdout {
	//	logrus.Out = logFile
	//} else {
	//	logrus.Out = multiLogger
	//}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func init() {
	usrDefTunFile := flag.Bool("udtuns", false, "List one or many user defined tunnel config files.  Files must be placed in the configuration directory.")
	configsDirFlag := flag.String("confDir", "configs/", "The directory for storing configuration files, also contains the instance configuration file.")
	logsDirFlag := flag.String("logsDir", "logs/", "The directory for storing logs.")
	templatesDirFlag := flag.String("templatesDir", "templates/", "The directory for storing templates.")

	flag.Parse()

	instanceConfig.ConfigDirectory = *configsDirFlag
	instanceConfig.LogDirectory = *logsDirFlag
	instanceConfig.TemplateDirectory = *templatesDirFlag

	configFileExists, configFileErr := file.Exists(filepath.Join(instanceConfig.ConfigDirectory + "ephi_instance.toml"))
	if !configFileExists || configFileErr != nil {
		logrus.WithFields(logrus.Fields{"Function": "init", "Message": "Failed to find instance configuration file.  Exiting application."}).Fatal(configFileErr)
	}

	regionsFileExists, regionsFileErr := file.Exists(filepath.Join(instanceConfig.ConfigDirectory + "regions.toml"))
	if !regionsFileExists || regionsFileErr != nil {
		logrus.WithFields(logrus.Fields{"Function": "init", "Message": "Failed to find regions configuration file.  Exiting application."}).Fatal(regionsFileErr)
	}

	logsDirExists, logsDirErr := file.Exists(instanceConfig.LogDirectory)
	if !logsDirExists {
		os.Mkdir(instanceConfig.LogDirectory, 0700)
	}
	if logsDirErr != nil {
		logrus.WithFields(logrus.Fields{"Function": "init", "Message": "Failed to find logs directory."}).Error(logsDirErr)
	}
	templatesDirExists, templatesDirErr := file.Exists(instanceConfig.TemplateDirectory)
	if !templatesDirExists {
		os.Mkdir(instanceConfig.TemplateDirectory, 0700)
	}
	if templatesDirErr != nil {
		logrus.WithFields(logrus.Fields{"Function": "init", "Message": "Failed to find templates directory."}).Error(templatesDirErr)
	}

	instanceConfig.loadInstanceConfig()
	if *configsDirFlag != "configs/" {
		if instanceConfig.ConfigDirectory == "configs/" {
			instanceConfig.ConfigDirectory = *configsDirFlag
		}
	}
	if *logsDirFlag != "logs/" {
		if instanceConfig.LogDirectory == "logs/" {
			instanceConfig.LogDirectory = *logsDirFlag
		}
	}
	if *templatesDirFlag != "templates/" {
		if instanceConfig.TemplateDirectory == "templates/" {
			instanceConfig.TemplateDirectory = *templatesDirFlag
		}
	}

	regionOptions.loadRegionOptions(instanceConfig.ConfigDirectory, "regions.toml")
	tunnelOptions.loadTunnelOptions(instanceConfig.ConfigDirectory, "ephi_default_tuns.toml")

	if *usrDefTunFile {
		instanceConfig.UserDefinedConfigs = flag.Args()
		logrus.WithFields(logrus.Fields{"Function": "main", "Configuration Files": instanceConfig.UserDefinedConfigs}).Info("Capturing user defined tunnel configuration files.")

		udTunnelOptions := TunnelOptions{}
		for _, udFile := range instanceConfig.UserDefinedConfigs {
			udFileExists, udFileErr := file.Exists(filepath.Join(instanceConfig.ConfigDirectory + udFile + ".toml"))
			if !udFileExists {
				logrus.WithFields(logrus.Fields{"Function": "init", "Message": udFile + ".toml configuration file not found."}).Error(udFileErr)
			} else {
				udTunnelOptions.loadTunnelOptions(instanceConfig.ConfigDirectory, udFile+".toml")
				for _, tunConfig := range udTunnelOptions.TunnelOptions {
					tunnelOptions.TunnelOptions = append(tunnelOptions.TunnelOptions, tunConfig)
				}
			}
		}
	}

	initDebugLog()
	instanceConfig.writeInstanceConfig()
}
