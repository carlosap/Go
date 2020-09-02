package cmd

import (
	c "github.com/Go/azuremonitor/config"
	"sync"
	"time"
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

)


