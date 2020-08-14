package dbcontext

import (
	"fmt"
	"testing"
)

func TestStorageaccount(t *testing.T) {

	t.Run("Test 1- Verify we can fetch all storageaccount access", func(t *testing.T) {
		sa := &Storageaccount{}
		sas, err := sa.GetAll()
		if err != nil {
			t.Errorf("error: %+v get all storageaccount accesses", err)
		}

		for i := 0; i < len(sas); i++ {
			fmt.Printf("Results: found a storageaccount  %d ]\n", sas[i].Resourceid)
		}

	})

}
