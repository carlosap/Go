package cmd

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Go/azuremonitor/db/cache"
	externalip "github.com/glendc/go-external-ip"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"unicode/utf8"
)

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

// returns internal ip and public ip
// you can also get this information https://myexternalip.com/raw
func getIP() ([]string, error) {
	var ips []string
	extIp := externalip.DefaultConsensus(nil, nil)
	ipTemp, _ := extIp.ExternalIP()
	if len(ipTemp.String()) > 0 {
		ips = append(ips, ipTemp.String())
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ips = append(ips, ip.String())
				//fmt.Printf("ip net: %s\n", ip.String())
			case *net.IPAddr:
				ip = v.IP
				//fmt.Printf("ip address: %s\n", ip.String())
				ips = append(ips, ip.String())
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ips, nil
		}
	}
	return ips, errors.New("no network connection detected")
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

//EncodeStringToBase32Value encodes string to base32
func EncodeStringToBase32Value(v string) string {
	d := []byte(v)
	return base32.StdEncoding.EncodeToString(d)
}

//DecodeBase32ToString takes a based32 string and returns string
func DecodeBase32ToString(v string) string {
	d, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return string(d)
}
