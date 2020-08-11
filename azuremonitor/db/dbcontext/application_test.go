package dbcontext

import (
	"fmt"
	"testing"
)

func TestApplication(t *testing.T) {

	t.Run("Test 1- Verify we can fetch all user access", func(t *testing.T) {
		app := &Application{}
		apps, err := app.GetAll()
		if err != nil {
			t.Errorf("error: %+v get all users accesses", err)
		}

		for i := 0; i < len(apps); i++ {
			fmt.Printf("Results: found a application  %d ]\n", apps[i].Applicationid)
		}

	})

}