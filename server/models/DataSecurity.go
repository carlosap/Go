package models

//DataSecurity makes an object that implements the PermissionCheckInterface
type DataSecurity struct {
	Groups                  `json:"groups,omitempty" cql:"groups"`
	Classification 			`json:"classification" cql:"classification"`
}

//Groups List
type Groups []string