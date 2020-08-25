package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/db/cache"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"unicode/utf8"
)

func IfErrorsPrintThem(errors []string) {
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
	}
}

func stringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f, nil
	}
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError {
		return 0, err
	}
	symbol := s[len(s)-size : len(s)]
	factor, ok := siFactors[symbol]
	if !ok {
		return 0, err
	}
	f, e := strconv.ParseFloat(s[:len(s)-len(symbol)], 64)
	if e != nil {
		return 0, err
	}
	return f * factor, nil
}

func clearTerminal() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("wrong platform")
	}
}

func clearCache(fileName string) {
	path := filepath.Join("cache", fileName)
	_ = os.Remove(path)
}

// Save saves a representation of v to the cachefolder
func saveCache(key string, v interface{}) error {
	c := &cache.Cache{}
	fileKey := guuid.New()
	path := filepath.Join("cache", fileKey.String())
	err := Save(path, v)
	if err != nil {
		return fmt.Errorf("failed to marshal ip information %v\n", err)
	}
	c.Set(key, fileKey.String())
	return err
}

func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	//b, err := json.MarshalIndent(v, "", "\t")
	b, err := json.MarshalIndent(v, "", "")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func LoadFromCache(cKey string, v interface{}) error {
	c := &cache.Cache{}
	cHashVal := c.Get(cKey)
	path := filepath.Join("cache", cHashVal)
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		c.Delete(cKey)
		return err
	}

	defer f.Close()
	return Unmarshal(f, v)
}

func loadFile(path string) ([]byte, error) {

	lock.Lock()
	defer lock.Unlock()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	return data, nil
}

func getCpuParallelCapabilities() (int, int) {
	var parallel int
	cpus := runtime.NumCPU()
	if cpus < 2 {
		parallel = 1
	} else {
		parallel = cpus - 1
	}

	if runtime.GOOS == "solaris" {
		parallel = 3
	}
	return parallel, cpus
}

func getStructNameByInterface(v interface{}) string {
	rv := reflect.ValueOf(v)
	typ := rv.Type()
	return typ.Name()
}

func saveCSV(filepath string, matrix [][]string) {

	if len(matrix) > 0 {

		if IsPathExists(filepath) == false {
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
