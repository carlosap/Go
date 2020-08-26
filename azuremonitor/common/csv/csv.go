package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/Go/azuremonitor/common/filesystem"
	"os"
)

func SaveMatrixToFile(filepath string, matrix [][]string) {

	if len(matrix) > 0 {

		if filesystem.IsPathExists(filepath) == false {
			_, err := os.Create(filepath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "--> %s\n", err)
				return
			}
		}

		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		w := csv.NewWriter(f)
		err = w.WriteAll(matrix)
		if err != nil {
			_ = f.Close()
			_, _ = fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		w.Flush()

		err = w.Error()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		_ = f.Close()
	}
}
