package costmanagement

import "time"

type BatchResponse struct {
	Responses []Responses `json:"responses"`
}
type Name struct {
	Value          string `json:"value"`
	LocalizedValue string `json:"localizedValue"`
}
type Data struct {
	TimeStamp string `json:"timeStamp"`
	Total float64 `json:"total"`
	Average   float64   `json:"average"`
}
type Timeseries struct {
	Metadatavalues []interface{} `json:"metadatavalues"`
	Data           []Data        `json:"data"`
}
type Value struct {
	ID                 string       `json:"id"`
	Type               string       `json:"type"`
	Name               Name         `json:"name"`
	DisplayDescription string       `json:"displayDescription"`
	Unit               string       `json:"unit"`
	Timeseries         []Timeseries `json:"timeseries"`
	ErrorCode          string       `json:"errorCode"`
}
type Content struct {
	Cost           int       `json:"cost"`
	Timespan       string `json:"timespan"`
	Interval       string    `json:"interval"`
	Value          []Value   `json:"value"`
	Namespace      string    `json:"namespace"`
	Resourceregion string    `json:"resourceregion"`
	Properties	   Properties `json:"properties"`

	Name       string      `json:"name"`
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Location   string      `json:"location"`
}
type Responses struct {
	HTTPStatusCode int     `json:"httpStatusCode"`
	Content        Content `json:"content"`
}

type HardwareProfile struct {
	VMSize string `json:"vmSize"`
}
type ImageReference struct {
	Publisher    string `json:"publisher"`
	Offer        string `json:"offer"`
	Sku          string `json:"sku"`
	Version      string `json:"version"`
	ExactVersion string `json:"exactVersion"`
}
type ManagedDisk struct {
	StorageAccountType string `json:"storageAccountType"`
	ID                 string `json:"id"`
}
type OsDisk struct {
	OsType       string      `json:"osType"`
	Name         string      `json:"name"`
	CreateOption string      `json:"createOption"`
	Caching      string      `json:"caching"`
	ManagedDisk  ManagedDisk `json:"managedDisk"`
	DiskSizeGB   int         `json:"diskSizeGB"`
}
type StorageProfile struct {
	ImageReference ImageReference `json:"imageReference"`
	OsDisk         OsDisk         `json:"osDisk"`
	DataDisks      []interface{}  `json:"dataDisks"`
}
type PatchSettings struct {
	PatchMode string `json:"patchMode"`
}
type WindowsConfiguration struct {
	ProvisionVMAgent       bool          `json:"provisionVMAgent"`
	EnableAutomaticUpdates bool          `json:"enableAutomaticUpdates"`
	PatchSettings          PatchSettings `json:"patchSettings"`
}
type OsProfile struct {
	ComputerName                string               `json:"computerName"`
	AdminUsername               string               `json:"adminUsername"`
	WindowsConfiguration        WindowsConfiguration `json:"windowsConfiguration"`
	Secrets                     []interface{}        `json:"secrets"`
	AllowExtensionOperations    bool                 `json:"allowExtensionOperations"`
	RequireGuestProvisionSignal bool                 `json:"requireGuestProvisionSignal"`
}
type NetworkInterfaces struct {
	ID string `json:"id"`
}
type NetworkProfile struct {
	NetworkInterfaces []NetworkInterfaces `json:"networkInterfaces"`
}
type Statuses struct {
	Code          string    `json:"code"`
	Level         string    `json:"level"`
	DisplayStatus string    `json:"displayStatus"`
	Message       string    `json:"message"`
	Time          time.Time `json:"time"`
}
type Status struct {
	Code          string `json:"code"`
	Level         string `json:"level"`
	DisplayStatus string `json:"displayStatus"`
	Message       string `json:"message"`
}
type ExtensionHandlers struct {
	Type               string `json:"type"`
	TypeHandlerVersion string `json:"typeHandlerVersion"`
	Status             Status `json:"status,omitempty"`
}
type VMAgent struct {
	VMAgentVersion    string              `json:"vmAgentVersion"`
	Statuses          []Statuses          `json:"statuses"`
	ExtensionHandlers []ExtensionHandlers `json:"extensionHandlers"`
}
type Disks struct {
	Name     string     `json:"name"`
	Statuses []Statuses `json:"statuses"`
}
type Extensions struct {
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	TypeHandlerVersion string     `json:"typeHandlerVersion"`
	Statuses           []Statuses `json:"statuses"`
}
type InstanceView struct {
	ComputerName     string       `json:"computerName"`
	OsName           string       `json:"osName"`
	OsVersion        string       `json:"osVersion"`
	VMAgent          VMAgent      `json:"vmAgent"`
	Disks            []Disks      `json:"disks"`
	Extensions       []Extensions `json:"extensions"`
	HyperVGeneration string       `json:"hyperVGeneration"`
	Statuses         []Statuses   `json:"statuses"`
}
type Properties struct {
	VMID              string          `json:"vmId"`
	HardwareProfile   HardwareProfile `json:"hardwareProfile"`
	StorageProfile    StorageProfile  `json:"storageProfile"`
	OsProfile         OsProfile       `json:"osProfile"`
	NetworkProfile    NetworkProfile  `json:"networkProfile"`
	ProvisioningState string          `json:"provisioningState"`
	InstanceView      InstanceView    `json:"instanceView"`
	AutoUpgradeMinorVersion bool     `json:"autoUpgradeMinorVersion"`
	Publisher               string   `json:"publisher"`
	Type                    string   `json:"type"`
	TypeHandlerVersion      string   `json:"typeHandlerVersion"`
}






