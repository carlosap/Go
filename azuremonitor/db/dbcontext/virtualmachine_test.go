package dbcontext

import (
	"fmt"
	"testing"
)

func TestVirtualmachine(t *testing.T) {

	t.Run("Test 1- Verify we can fetch all virtualmachine access", func(t *testing.T) {
		vm := &Virtualmachine{}
		vms, err := vm.GetAll()
		if err != nil {
			t.Errorf("error: %+v get all virtualmachine accesses", err)
		}

		for i := 0; i < len(vms); i++ {
			fmt.Printf("Results: found a virtualmachine  %d ]\n", vms[i].Resourceid)
		}

	})

}
