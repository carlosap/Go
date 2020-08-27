package cmd

import (
	"bytes"
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
)


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

func getStructNameByInterface(v interface{}) string {
	rv := reflect.ValueOf(v)
	typ := rv.Type()
	return typ.Name()
}
