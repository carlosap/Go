package dbcontext

import (
	"testing"
)

func TestDbContext(t *testing.T) {

	t.Run("Test 1- Verify we can create new dbcontext object", func(t *testing.T) {
		db, err := NewDbContext()
		if err != nil {
			t.Fatalf("error: failed to connect to db %v", err)
		}

		_ = db.Close()
	})

}
