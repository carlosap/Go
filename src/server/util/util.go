package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// List of functions to call before exiting
var (
	exitFunctions    = make([]func(), 0)
	exitFunctionLock = sync.RWMutex{}
	exitFunctionOnce = sync.Once{}
)

func init() {
	registerForSignal()
}

// Fibonacci computes the Nth number in the fibonacci sequence.
func Fibonacci(n int) int {
	if n <= 1 {
		return 1
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// GetSetEnv Gets and potentially sets environment variable to the fallback value.
// Returns environment-set value if present, fallback otherwise
func GetSetEnv(envVar, fallback string) string {
	envString := os.Getenv(envVar)
	if envString == "" {
		envString = fallback
		if err := os.Setenv(envVar, envString); err != nil {
			log.Printf("Unable to set env %s = %s: %v", envVar, envString, err)
		}
	}

	return envString
}

// Register for the common process killing signals
func registerForSignal() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		_ = <-sig
		exitFunctionLock.RLock()
		if len(exitFunctions) > 0 {
			log.Printf("Catching exit signal...")
			RunExitSignals()
		}
		exitFunctionLock.RUnlock()
		os.Exit(1)
	}()
}

// ExitOnOrphan exits the program when it becomes orphaned.
// When a process is orphaned, its parent becomes init, which is always PID 1.
func ExitOnOrphan() {
	for range time.Tick(time.Minute) {
		if os.Getppid() == 1 {
			RunExitSignals()
			log.Fatal("Exiting because process was orphaned.")
		}
	}
}

// RunExitSignals fires each registered exit function when an exit signal is
// received.
func RunExitSignals() {
	exitFunctionOnce.Do(runExitSignals)
}

func runExitSignals() {
	exitFunctionLock.RLock()
	for _, f := range exitFunctions {
		f()
	}
	exitFunctionLock.RUnlock()
}

// CatchExitSignal is an external interface to register your exit function.
func CatchExitSignal(f func()) {
	exitFunctionLock.Lock()
	exitFunctions = append(exitFunctions, f)
	exitFunctionLock.Unlock()
}

// Max returns the greater of two int64s.
func Max(value1, value2 int64) int64 {
	if value1 > value2 {
		return value1
	}
	return value2
}

// MinInt returns the minimum of two ints.
func MinInt(value1, value2 int) int {
	if value1 < value2 {
		return value1
	}
	return value2
}

// MaxInt returns the maximum of two ints.
func MaxInt(value1, value2 int) int {
	if value1 > value2 {
		return value1
	}
	return value2
}

// Concat appends one byte slice to the end of the other and returns the result.
func Concat(old1, old2 []byte) []byte {
	return append(old1, old2...)
}

// SpliceFromStringSlice removes a string from a slice of strings.
func SpliceFromStringSlice(list []string, index int) []string {
	newLen := len(list) - 1

	// Shift the list and set the last item to an empty string to
	// remove the reference, then truncate.
	copy(list[index:], list[index+1:])
	list[newLen] = ""
	list = list[:newLen]

	return list
}

// IndexOfString Searches for 'searchString' in 'list' and returns the index for
// the first match or a -1 if the string is not found.
// Note:  This is not an efficient function when the list is large.
func IndexOfString(list []string, searchString string) int {
	foundIdx := -1
	for idx, str := range list {
		if str == searchString {
			foundIdx = idx
			break
		}
	}
	return foundIdx
}

// IndexOfInt searches for 'num' in 'list' and returns the index for the first
// match or a -1 if the integer is not found.
// Note: This is not an effective function when the list is large.
func IndexOfInt(list []int, num int) int {
	for i, x := range list {
		if x == num {
			return i
		}
	}
	return -1
}

// StringSlicesEqual returns true if the two string slices have the same content.
func StringSlicesEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Create a map for all the strings in s2.
	s2Map := make(map[string]struct{}, len(s2))
	for _, s := range s2 {
		s2Map[s] = struct{}{}
	}

	for _, s := range s1 {
		if _, ok := s2Map[s]; !ok {
			return false
		}
	}
	return true
}

//IsSliceSubset returns true if all elements in s1 are found in s2
func IsSliceSubset(s1, s2 []string) bool {
	for _, s := range s1 {
		found := false
		for _, sub2 := range s2 {
			if s == sub2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// StringSliceMapDiffs takes two string arrays, s1 and s2, and returns all
// elements of s1 that are not in s2, up to 'max' elements.  The strings
// will be returned as keys to a map.
func StringSliceMapDiffs(s1, s2 []string, max int) map[string]struct{} {
	// Create a map for all the strings in s2.
	s2Map := make(map[string]struct{}, len(s2))
	for _, s := range s2 {
		s2Map[s] = struct{}{}
	}

	// If an element of s1 isn't found in the s2Map, save it in the
	// returned map (diffs), unless we've reached the maximum length
	// of differences that we're allowed to return.
	diffs := make(map[string]struct{}, max)
	size := 0
	for _, s := range s1 {
		if _, ok := s2Map[s]; !ok {
			if size < max {
				diffs[s] = struct{}{}
				size++
			} else {
				break
			}
		}
	}
	return diffs
}

// ShouldZipResponse adds the gzip compression scheme token to the HTTP request.
func ShouldZipResponse(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

// LogHTTPRequest logs the function name and remoteAddress of a HTTP request.
func LogHTTPRequest(funcName, remoteAddress string) {
	log.Printf("Got %s request from %s\n", funcName, remoteAddress)
}

// PipeStruct specifies a slice of bytes; used to represent piped stdin data.
type PipeStruct struct {
	Data []byte `json:"data"`
}

// ReadFromStdin reads from stdin until requested to stop.
func ReadFromStdin(stop <-chan struct{}) <-chan []byte {
	outChan := make(chan []byte)
	go func() {
		jReader := json.NewDecoder(os.Stdin)
		p := PipeStruct{}
		var err error
		for err == nil {
			err = jReader.Decode(&p)
			select {
			case outChan <- p.Data:
			case <-stop:
				return
			}
		}
		if err != nil {
			log.Printf("Error reading from stdin: %+v", err)
		}
		close(outChan)
	}()
	return outChan
}

// StdinReader is a collection consisting of a JSON decoder and slice of bytes.
type StdinReader struct {
	jReader *json.Decoder
	buf     []byte
}

// NewStdinReader creates and initializes a StdinReader struct.
func NewStdinReader() StdinReader {
	s := StdinReader{}
	s.buf = make([]byte, 0)
	s.jReader = json.NewDecoder(os.Stdin)
	return s
}

// Read decodes incoming JSON data and returns the number of bytes read.
func (s StdinReader) Read(p []byte) (n int, err error) {
	for len(s.buf) < len(p) && err == nil {
		// Grab the next packet and append the bytes to our buffer
		ps := PipeStruct{}
		err = s.jReader.Decode(&ps)
		s.buf = append(s.buf, ps.Data...)
	}

	n = copy(p, s.buf)
	s.buf = s.buf[n:]
	return
}

// WriteToStdout encodes incoming data and pipes it to stdout.
func WriteToStdout(outChan <-chan []byte) error {
	p := PipeStruct{}
	jWriter := json.NewEncoder(os.Stdout)
	var d []byte
	var i int
	for d = range outChan {
		p.Data = d
		err := jWriter.Encode(p)
		if err != nil {
			return fmt.Errorf("Error writing to stdout: %+v", err)
		}
		i++
	}
	return nil
}

// StdoutWriter represents an JSON encoder.
type StdoutWriter struct {
	jWriter *json.Encoder
}

// NewStdoutWriter creates and initializes a StdoutWriter struct.
func NewStdoutWriter() StdoutWriter {
	return StdoutWriter{jWriter: json.NewEncoder(os.Stdout)}
}

// Write encodes outgoing JSON data and returns the number of bytes written.
func (s StdoutWriter) Write(p []byte) (int, error) {
	n := len(p)
	return n, s.jWriter.Encode(PipeStruct{Data: p})
}

// RunCmd runs the given executable with the given arguments.
//  executable: executable name(or full path if its not in the path)
//  args: list of arguments as strings
//  timeout: if < 0 wait for the commmand to return before exiting
//           if == 0 return  immediately
//           if > 0 max number of seconds process will run before killing process
//  verbose: whether to print information about how the call goes
func RunCmd(executable string, args []string, timeout int, verbose bool) (*exec.Cmd, bytes.Buffer, bytes.Buffer) {
	cmd := exec.Command(executable, args...)

	if verbose {
		log.Print(cmd)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	var err error
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if timeout < 0 {
		err = cmd.Run()
	} else {
		err = cmd.Start()
		//running in timeout mode
		if err == nil && timeout > 0 {
			//kill running process if timeout expires
			kill := time.AfterFunc(time.Second*time.Duration(timeout), func() { cmd.Process.Kill() })

			//wait until process returns,
			//process will return if timeout expires and kill signal is sent
			err = cmd.Wait()

			//stop time if it's still running,
			//we don't want to send a signal to a process in the future with a recovered pid
			kill.Stop()
		}
	}

	if verbose {
		if err != nil {
			fmt.Println("Error running ", executable,
				fmt.Sprint(err)+": "+stderr.String())
		} else {
			fmt.Println("No Error: " + stderr.String())
		}
		fmt.Println("Result: " + out.String())
	}

	return cmd, out, stderr
}

// There are leap second offsets from 1972 onwards, but we're only going to
// consider from 2006 and onwards.  As time goes by, additional entries will
// be needed for the leapSecondTable.
//
// Further info on leap seconds: https://en.wikipedia.org/wiki/Leap_second
//
// ... "In 1972, the leap-second system was introduced so that the broadcast
// UTC seconds could be made exactly equal to the standard SI second, while
// still maintaining the UTC time of day and changes of UTC date synchronized
// with those of UT1 (the solar time standard that superseded GMT).[8] By
// then, the UTC clock was already 10 seconds behind TAI, which had been
// synchronized with UT1 in 1958, but had been counting true SI seconds
// since then. After 1972, both clocks have been ticking in SI seconds,
// so the difference between their readouts at any time is 10 seconds
// plus the total number of leap seconds that have been applied to UTC
// (36 seconds in July 2015)."

type leapSecondInfo struct {
	julianDay int     // Julian day at time of leap second
	julianSec float64 // Julian secondsOfDay at time of leap second
	offset    float64 // leap seconds offset
}

var leapSecondTable = []leapSecondInfo{
	{2453736, 43233.0, 33.0}, // January 1, 2006 00:00:00 UTC
	{2454832, 43234.0, 34.0}, // January 1, 2009 00:00:00 UTC
	{2456109, 43235.0, 35.0}, // July 1, 2012 00:00:00 UTC
	{2457204, 43236.0, 36.0}, // July 1, 2015 00:00:00 UTC
	{2457749, 43237.0, 37.0}, // January 1, 2017 00:00:00 UTC
}

// Get the leap second offset given a Julian day.
func getLeapSecondOffset(julianDay int) float64 {
	lastIndex := len(leapSecondTable) - 1
	for i := lastIndex; i > 0; i-- {
		if julianDay >= leapSecondTable[i].julianDay {
			return leapSecondTable[i].offset
		}
	}

	if julianDay < leapSecondTable[0].julianDay {
		log.Printf("Error:  Julian day '%d' represents a date less than "+
			"January 1, 2006", julianDay)
		// Use the earliest offset we have ...
	}
	return leapSecondTable[0].offset
}

// GetJulianTime takes a time and returns the Julian date in days and the
// number of seconds into that day.
// This formula for the day calculation was taken from Rob Pike's suggestion
// found here:  https://groups.google.com/forum/#!topic/golang-nuts/biCftPuMbDk
func GetJulianTime(t time.Time) (int, float64) {
	// Julian date, in seconds, of the "Format" standard time, i.e.,
	// (1/2/2006 15:04:05 UTC).  (See:
	// http://www.onlineconversion.com/julian_date.htm)
	const julian = 2453738.4195

	// Easiest way to get the time.Time of the Unix time.
	// (See comments for the UnixDate in package Time.)
	unix := time.Unix(1136239445, 0)
	const oneDay = float64(86400. * time.Second)

	// Take Golang's standard time and add the days since then.
	day := int(julian + float64(t.Sub(unix))/oneDay)

	// Get time info for this date.
	hour, minute, second := t.Clock()
	var milliseconds = float64(t.Nanosecond()) / 1000000.0

	// Julian dates are noon-based, so adjust accordingly
	hour -= 12
	if hour < 0 {
		hour += 24
	}
	var secondsOfDay = float64(second+(hour*60*60)+(minute*60)) +
		(milliseconds * .001)

	// Add the leap second offset.
	secondsOfDay += getLeapSecondOffset(day)

	return day, secondsOfDay
}
