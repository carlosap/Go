package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

//ReadSeekCloser interface for holding file io handles
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

//GetFile Retrieves a file pointer and the files mtime
func GetFile(fileName string) (ReadSeekCloser, time.Time, error) {
	var modTime time.Time
	fi, err := os.Open(fileName)
	if err == nil {
		fInfo, _ := fi.Stat()
		if fInfo != nil {
			modTime = fInfo.ModTime()
		}
	}
	return fi, modTime, err
}

// GetModTime makes sure the file exists
func GetModTime(filePath string) (time.Time, error) {
	modTime := time.Time{}
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return modTime, err
	}
	modTime = fInfo.ModTime()
	return modTime, nil
}

//CreateZipArchive takes in a map of file names and their file paths, and a map of file names and the metadata
//it then zips these files and returns a reader
func CreateZipArchive(archiveName string, files map[string]string, metas map[string]string) (io.ReadSeeker, time.Time, error) {
	//Create zip file
	buf := new(bytes.Buffer)
	zipper := zip.NewWriter(buf)
	modTime := time.Time{}

	//Read and append each file
	for fileName, filePath := range files {
		ext := filepath.Ext(fileName)
		base := strings.Replace(fileName, ext, "", 1)
		data, err := ReadFile(filePath)
		if err != nil {
			return nil, modTime, fmt.Errorf("CreateZipArchive was unable to read data for file %s: %s", filePath, err.Error())
		}

		z, err := zipper.Create(fileName)
		if err != nil {
			return nil, modTime, fmt.Errorf("CreateZipArchive was unable to create zip for file %s: %s", filePath, err.Error())
		}
		_, err = z.Write(data)
		if err != nil {
			return nil, modTime, fmt.Errorf("CreateZipArchive was unable to write data for file %s to zip: %s", filePath, err.Error())
		}
		if metadata, ok := metas[filePath]; ok {
			metaName := fmt.Sprintf("%s%s", base, ".meta")
			metaContent := []byte(metadata)

			x, err := zipper.Create(metaName)
			if err != nil {
				return nil, modTime, fmt.Errorf("CreateZipArchive was unable to create zip for metadata %s: %s", filePath, err.Error())
			}
			_, err = x.Write(metaContent)
			if err != nil {
				return nil, modTime, fmt.Errorf("CreateZipArchive was unable to write metadata file for %s to zip: %s", filePath, err.Error())
			}
		}
	}
	//Cleanup
	err := zipper.Close()
	if err != nil {
		return nil, modTime, fmt.Errorf("Unable to close zip writer due to error: %s", err.Error())
	}

	return bytes.NewReader(buf.Bytes()), time.Now(), nil
}

//ReadFile will read a file from wherever it is stored and return the bytes
func ReadFile(filePath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return bytes, fmt.Errorf("ReadFile error while reading bytes from path %s: %s", filePath, err.Error())
	}
	return bytes, nil
}

//TrimExt removes the suffix from a file path
func TrimExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

//IsSameTimestamp tells whether the files are the same
//compares the timestamps of the two files.
func IsSameTimestamp(sourcePath, destPath string) (bool, error) {
	// Get the original modification time of the file.
	sourceTime, err := GetModTime(sourcePath)
	if err != nil {
		return false, err
	}
	// Get the thjmbnail modification time of the file.
	destTime, err := GetModTime(destPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	// Check if they are the same time
	// Implies that thes are the same file, and we can just use that.
	return sourceTime == destTime, nil
}

//RunCmdSafely runs the command and logs errors
func RunCmdSafely(cmdStr string, args []string, workdir string, env []string) error {
	// run the command and return err
	cmd := exec.Command(cmdStr, args...)
	if len(env) > 0 {
		cmd.Env = env
	}
	if workdir != "" {
		cmd.Dir = workdir
	}

	// Run the command and pipe errors to log
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	errPipeScanner := bufio.NewScanner(stderr)
	for errPipeScanner.Scan() {
		log.Println(errPipeScanner.Text())
	}
	if err := errPipeScanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}

// Taken from go1.6 tempFile.go, add suffix to TempFile()

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var tempRand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := tempRand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	tempRand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, and ending with suffix (ea file extension .html)
// opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file.  The caller can use f.Name()
// to find the pathname of the file.  It is the caller's responsibility
// to remove the file when no longer needed.
func TempFile(dir, prefix string, suffix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				tempRand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

//ReadYamlFile unmarshal's a yaml file into given object
func ReadYamlFile(fileName string, i interface{}) error {
	return UnmarshalFile(fileName, i, yaml.Unmarshal)
}

//UnmarshalFile A generic unmarshaler from a file
func UnmarshalFile(fileName string, i interface{}, f func([]byte, interface{}) error) error {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	return f(d, i)
}
