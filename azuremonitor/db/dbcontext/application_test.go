package dbcontext

import (
	"fmt"
	"testing"
)

func TestApplication(t *testing.T) {

	t.Run("Test 1- Verify we can fetch all application access", func(t *testing.T) {
		app := &Application{}
		apps, err := app.GetAll()
		if err != nil {
			t.Errorf("error: %+v get all application accesses", err)
		}

		for i := 0; i < len(apps); i++ {
			fmt.Printf("Results: found a application  %d ]\n", apps[i].Applicationid)
		}

	})

}

func Test(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"carlos"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name != "hermes"  {
				t.Errorf("failed to pass")
			}
		})
	}
}
