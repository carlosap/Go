package errors

import (
	"fmt"
	"os"
)

func IfErrorsPrintThem(errors []string) {
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
	}
}
