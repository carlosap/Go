package main

import (
	"path/filepath"

	"github.com/aagon00/POILibs/file"
	"github.com/sirupsen/logrus"
)

type TunnelConfig struct {
	Name            string       `toml:"name" json:"name"`
	Username        string       `toml:"username" json:"username"`
	PointOfPresence int          `toml:"point_of_presence" json:"point_of_presence"`
	PointOfEntry    int          `toml:"point_of_entry" json:"point_of_entry"`
	Nodes           []NodeConfig `toml:"nodes" json:"nodes"`
}

type NodeConfig struct {
	CloudProvider        string      `toml:"cloud_provider" json:"cloud_provider"`
	Region               string      `toml:"region" json:"region"`
	ImageName            string      `toml:"image_name" json:"image_name"`
	PublicIP             string      `toml:"public_ip" json:"public_ip"`
	PrivateIP            string      `toml:"private_ip" json:"private_ip"`
	PrivateSubnet        string      `toml:"private_subnet" json:"private_subnet"`
	VirtualSubnet        string      `toml:"virtual_subnet" json:"virtual_subnet"`
	SSHIPs               []string    `toml:"ssh_ips" json:"ssh_ips"`
	SSHPort              string      `toml:"ssh_port" json:"ssh_port"`
	VPNIPs               []string    `toml:"vpn_ips" json:"vpn_ips"`
	VPNPorts             []string    `toml:"vpn_ports" json:"vpn_ports"`
	WebIPs               []string    `toml:"web_ips" json:"web_ips"`
	WebPorts             []string    `toml:"web_ports" json:"web_ports"`
	AllowAllTrafficOut   bool        `toml:"allow_all_traffic_out" json:"allow_all_traffic_out"`
	IPV6AcceptTemplate   string      `toml:"ipv6_accept_template" json:"ipv6_accept_template"`
	IPV6AcceptRules      []string    `toml:"ipv6_accept_rules" json:"ipv6_accept_rules"`
	IPV4AcceptTemplate   string      `toml:"ipv4_accept_template" json:"ipv4_accept_template"`
	IPV4AcceptRules      []string    `toml:"ipv4_accept_rules" json:"ipv4_accept_rules"`
	IPV6DropTemplate     string      `toml:"ipv6_drop_template" json:"ipv6_drop_template"`
	IPV6DropRules        []string    `toml:"ipv6_drop_rules" json:"ipv6_drop_rules"`
	IPV4DropTemplate     string      `toml:"ipv4_drop_template" json:"ipv4_drop_template"`
	IPV4DropRules        []string    `toml:"ipv4_drop_rules" json:"ipv4_drop_rules"`
	IPSecConfigTemplate  string      `toml:"ipsec_config_template" json:"ipsec_config_template"`
	IPSecConfig          []string    `toml:"ipsec_config" json:"ipsec_config"`
	IPSecSecretsTemplate string      `toml:"ipsec_secrets_template" json:"ipsec_secrets_template"`
	IPSecSecrets         []string    `toml:"ipsec_secrets" json:"ipsec_secrets"`
	CloudConfig          interface{} `toml:"cloud_config" json:"cloud_config"`
}

type AWSConfig struct {
	VPCID                 string `toml:"vpc_id" json:"vpc_id"`
	SubnetID              string `toml:"subnet_id" json:"subnet_id"`
	RouteTableID          string `toml:"route_table_id" json:"route_table_id"`
	InternetGatewayID     string `toml:"internet_gateway_id" json:"internet_gateway_id"`
	InstanceID            string `toml:"instance_id" json:"instance_id"`
	VolumeID              string `toml:"volume_id" json:"volume_id"`
	SecurityGroupID       string `toml:"security_group_id" json:"security_group_id"`
	ElasticIPAllocationID string `toml:"elastic_ip_allocation_id" json:"elastic_ip_allocation_id"`
	KeyPairName           string `toml:"key_pair_name" json:"key_pair_name"`
}

type AzureConfig struct {
	ComputerName         string `toml:"computer_name" json:"computer_name"`
	NetworkSecurityGroup string `toml:"network_security_group" json:"network_security_group"`
	VirtualNetwork       string `toml:"virtual_network" json:"virtual_network"`
	DiskName             string `toml:"disk_name" json:"disk_name"`
	PublicIPAddress      string `toml:"public_ip_address" json:"public_ip_address"`
	NetworkInterface     string `toml:"network_interface" json:"network_interface"`
	StorageAccount       string `toml:"storage_account" json:"storage_account"`
}

type GCEConfig struct {
	VPCNetworkName      string   `toml:"vpc_network_name" json:"vpc_network_name"`
	StaticIPAddressName string   `toml:"static_ip_address_name" json:"static_ip_address_name"`
	FirewallRuleName    []string `toml:"firewall_rule_name" json:"firewall_rule_name"`
	VMInstanceName      string   `toml:"vm_instance_name" json:"vm_instance_name"`
	DiskName            string   `toml:"disk_name" json:"disk_name"`
}

type TunnelOptions struct {
	TunnelOptions []TunnelConfig `toml:"tunnel_options" json:"tunnel_options"`
}

// loadTunnelOptions loads the TOML file for tunnel configuration settings. Must pass file name including extension '.toml'.
func (to *TunnelOptions) loadTunnelOptions(configDirectory, fileName string) {
	err := file.SafeReadTomlFile(filepath.Join(configDirectory+fileName), &to)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "loadTunnelOptions", "Configuration File": fileName, "Message": "Error loading toml configuration file."}).Fatal(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "loadTunnelOptions", "Configuration File": fileName}).Info("Successfully loaded toml file.")
}

// writeTunnelOptions writes the TOML file for tunnel configuration settings. Must pass file name including extension '.toml'.
func (to *TunnelOptions) writeTunnelOptions(configDirectory, fileName string) {
	err := file.SafeWriteTomlFile(filepath.Join(configDirectory+fileName), to)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "writeTunnelOptions", "Configuration File": fileName, "Message": "Error in writing toml contents."}).Error(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "writeTunnelOptions", "Configuration File": fileName}).Info("Successfully wrote toml file.")
}

type RegionZone struct {
	Provider          string   `toml:"provider" json:"provider"`
	RegionName        string   `toml:"region_name" json:"region_name"`
	RegionDescription string   `toml:"region_description" json:"region_description"`
	Location          string   `toml:"location" json:"location"`
	NodeType          string   `toml:"node_type" json:"node_type"`
	Active            bool     `toml:"active" json:"active"`
	Zones             []string `toml:"zones" json:"zones"`
}

type RegionZonesOptions struct {
	RegionsAndZones []RegionZone `toml:"regions_and_zones" json:"regions_and_zones"`
}

// loadRegionOptions loads the TOML file for region configuration settings. Must pass file name including extension '.toml'.
func (rzo *RegionZonesOptions) loadRegionOptions(configDirectory, fileName string) {
	err := file.SafeReadTomlFile(filepath.Join(configDirectory+fileName), &rzo)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "loadRegionOptions", "Configuration File": fileName, "Message": "Error loading toml configuration file."}).Fatal(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "loadRegionOptions", "Configuration File": fileName}).Info("Successfully loaded toml file.")
}

// writeRegionOptions writes the TOML file for region configuration settings. Must pass file name including extension '.toml'.
func (rzo *RegionZonesOptions) writeRegionOptions(configDirectory, fileName string) {
	err := file.SafeWriteTomlFile(filepath.Join(configDirectory+fileName), rzo)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "writeRegionOptions", "Configuration File": fileName, "Message": "Error in writing toml contents."}).Error(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "writeRegionOptions", "Configuration File": fileName}).Info("Successfully wrote toml file.")
}

type EphIConfig struct {
	Name               string   `toml:"name" json:"name"`
	LogToFile          bool     `toml:"log_to_file" json:"log_to_file"`
	LogToStdout        bool     `toml:"log_to_stdout" json:"log_to_stdout"`
	LogDirectory       string   `toml:"log_directory" json:"log_directory"`
	ConfigDirectory    string   `toml:"config_directory" json:"config_directory"`
	TemplateDirectory  string   `toml:"template_directory" json:"template_directory"`
	UserDefinedConfigs []string `toml:"user_defined_configs" json:"user_defined_configs"`
}

// loadInstanceConfig loads the TOML file for node configuration settings.
func (ic *EphIConfig) loadInstanceConfig() {
	err := file.SafeReadTomlFile(filepath.Join(ic.ConfigDirectory+"ephi_instance.toml"), &ic)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "loadInstanceConfig", "Configuration File": "ephi_instance.toml", "Message": "Error loading toml configuration file."}).Fatal(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "loadInstanceConfig", "Configuration File": "ephi_instance.toml"}).Info("Successfully loaded toml file.")
}

// writeInstanceConfig writes the TOML file for node configuration settings.
func (ic *EphIConfig) writeInstanceConfig() {
	err := file.SafeWriteTomlFile(filepath.Join(ic.ConfigDirectory+"ephi_instance.toml"), ic)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Function": "writeInstanceConfig", "Configuration File": "ephi_instance.toml", "Message": "Error in writing toml contents."}).Error(err)
	}

	logrus.WithFields(logrus.Fields{"Function": "writeInstanceConfig", "Configuration File": "ephi_instance.toml"}).Info("Successfully wrote toml file.")
}
