package filesystem

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func IsPathExists(path string) bool {
	_, err := os.Lstat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// error e.g. permission denied
	return false
}

func IsDirectoryExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Folder does not exist.")
		return false
	}
	return true
}
func RemoveFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("file doesn't exist\n")
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		return false
	}
	return true
}
func RemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("failed to remove directory ('%s') failed with '%s'\n", path, err)
	}
}
