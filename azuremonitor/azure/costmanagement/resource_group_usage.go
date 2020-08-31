package costmanagement

import (
	"fmt"
	"github.com/Go/azuremonitor/common/filesystem"
)

type ResourceGroupUsage struct {}


func (rgu *ResourceGroupUsage) RunAll() {

	// virtual machines
	rgu.virtualMachines()
}

func (rgu *ResourceGroupUsage) virtualMachines() {
	vm := VirtualMachine{}
	vm.ExecuteRequest(&vm)
	vm.Print()
	if SaveCsv {
		filesystem.RemoveFile(CsvRguReportName)
		vm.WriteCSV(CsvRguReportName)
		fmt.Printf("Done. report was generated - %s\n", CsvRguReportName)
	}
}
