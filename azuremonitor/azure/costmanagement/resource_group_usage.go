package costmanagement

import (
	"fmt"
	"github.com/Go/azuremonitor/common/filesystem"
)

type ResourceGroupUsage struct {}

var (
	csvStorageDiskReport = "storage_disk.csv"
	csvVirtualMachineReport = "virtual_machine.csv"
	csvLogicAppWorkflowReport = "logicapp_workflow.csv"
	csvStorageAccountReport = "storage_account.csv"
)

func (rgu *ResourceGroupUsage) RunAll() {

	//Three Node of Resources
	if len(Resources) <= 0 {
		rg := ResourceGroupCost{}
		rg.ExecuteRequest(&rg)
	}

	// virtual machines
	rgu.virtualMachines()

	// storage disk - vm
	rgu.storageDisk()

	// logicapps workflow
	rgu.logicAppWorkflows()

	// storage storageaccounts (files, blog, queue)
	rgu.storageStorageAccount()


}

func (rgu *ResourceGroupUsage) virtualMachines() {
	vm := VirtualMachine{}
	vm.ExecuteRequest(&vm)
	vm.Print()
	if SaveCsv {
		filesystem.RemoveFile(csvVirtualMachineReport)
		vm.WriteCSV(csvVirtualMachineReport)
		fmt.Printf("Done. report was generated - %s\n\n\n\n\n", csvVirtualMachineReport)
	}
}

func (rgu *ResourceGroupUsage) storageDisk() {
	sd := StorageDisk{}
	sd.ExecuteRequest(&sd)
	sd.Print()
	if SaveCsv {
		filesystem.RemoveFile(csvStorageDiskReport)
		sd.WriteCSV(csvStorageDiskReport)
		fmt.Printf("Done. report was generated - %s\n\n\n\n\n", csvStorageDiskReport)
	}
}
func (rgu *ResourceGroupUsage) storageStorageAccount() {
	st := StorageAccount{}
	st.ExecuteRequest(&st)
	st.Print()
	if SaveCsv {
		filesystem.RemoveFile(csvStorageAccountReport)
		st.WriteCSV(csvStorageAccountReport)
		fmt.Printf("Done. report was generated - %s\n\n\n\n\n", csvStorageAccountReport)
	}
}

func (rgu *ResourceGroupUsage) logicAppWorkflows() {
	lg := LogicAppWorkFlow{}
	lg.ExecuteRequest(&lg)
	lg.Print()
	if SaveCsv {
		filesystem.RemoveFile(csvLogicAppWorkflowReport)
		lg.WriteCSV(csvLogicAppWorkflowReport)
		fmt.Printf("Done. report was generated - %s\n\n\n\n\n", csvLogicAppWorkflowReport)
	}
}
