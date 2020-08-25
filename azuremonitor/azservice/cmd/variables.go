package cmd

import (
	c "github.com/Go/azuremonitor/config"
	"sync"
	"time"
)

const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000
)

var (
	configuration    c.CmdConfig
	layoutISO        = "2006-01-02"
	startDate        string
	endDate          string
	saveDb           bool
	saveCsv          bool
	ignoreZeroCost   bool
	ctr              = 0
	startTime        time.Time
	lock             sync.Mutex
	developer        string
	version          = "0.3"
	parallel, cpus   = getCpuParallelCapabilities()
	Methods          = &RequestMethods{POST: "POST", GET: "GET"}
	csvRgcReportName = "resource_group_cost.csv"
	csvRguReportName = "resource_group_usage.csv"
)

var siFactors = map[string]float64{
	"":  1e0,
	"k": 1e3,
	"M": 1e6,
	"G": 1e9,
	"T": 1e12,
	"P": 1e15,
	"E": 1e18,
	"Z": 1e21,
	"Y": 1e24,
	"K": 1e3,
	"B": 1e9,
}
