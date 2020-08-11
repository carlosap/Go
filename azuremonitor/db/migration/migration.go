package main

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Go/azuremonitor/config"
	"github.com/fatih/color"
	"github.com/lib/pq"
)

// Direction Builds the operations that will migrate the database 'up/down'.
type Direction int

// File system.io files struct
type File struct {
	Path      string
	FileName  string
	Version   uint64
	Name      string
	Content   []byte
	Direction Direction
}

// Files slice of file
type Files []File

// Migration Files
type MFile struct {
	Version  uint64
	UpFile   *File
	DownFile *File
}

// MFiles migration files
type MFiles []MFile

// DB handler to exec SQL commands
type DBDriver struct {
	db     *sql.DB
	tx     *sql.Tx
	Schema string
}

// Driver sql interface to allow custum implemenation
type Driver interface {
	Initialize() error
	Close() error
	Migrate(file File, pipe chan interface{})
	Version() (uint64, error)
	ResetVersion() (uint64, error)
}

var driversMu sync.Mutex
var drivers = make(map[string]Driver)

// TODO::Issue #254- pgmigrator get connectionstring from config .toml
var Scheme = "postgres"

const Schema = "azuremonitor"

//var connectionString = "host=localhost port=5432 user=postgres password=password dbname=elysium sslmode=disable"
var interrupts = true
var fileExpression = `^([0-9]+)_(.*)\.(up|down)\.%s$`
var fileExtension = "sql"

// flags are commands supported.
var version = flag.Bool("version", false, "Show version")
var filepath = flag.String("filepath", "", "")

const (
	// Migration Direction Up
	UpDirection Direction = +1
	// Migration Direction Down
	DownDirection Direction = -1

	Version   = "1.0.0"
	tableName = "schema_migrations"
)

// Sorts
func (mf MFiles) Len() int {
	return len(mf)
}

func (mf MFiles) Less(i, j int) bool {
	return mf[i].Version < mf[j].Version
}

func (mf MFiles) Swap(i, j int) {
	mf[i], mf[j] = mf[j], mf[i]
}

func init() {
	RegisterDriver("postgres", &DBDriver{})
}

func main() {
	flag.Usage = func() {
		PrintUsage()
	}
	flag.Parse()

	cmd := flag.Arg(0)
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *filepath == "" {
		*filepath, _ = os.Getwd()
	}

	fmt.Printf("selected command: %s\n", cmd)
	fmt.Printf("parse directory file path: %s\n", *filepath)

	switch cmd {

	case "up":
		pipe := NewPipeChannel()
		go RunMigrationUpToDate(pipe, *filepath)
		ok := writePipe(pipe)
		if !ok {
			os.Exit(1)
		}

	case "reset":
		pipe := NewPipeChannel()
		fmt.Printf("Resetting Migration Version to initial state\n")
		if err := Reset(); err != nil {
			fmt.Printf("resetting error: %v\n", err)
			os.Exit(1)
		}

		go RunMigrationUpToDate(pipe, *filepath)
		ok := writePipe(pipe)
		if !ok {
			os.Exit(1)
		}
	default:
		PrintUsage()
		os.Exit(1)
	case "help":
		PrintUsage()

	}
}

func PrintUsage() {
	os.Stderr.WriteString(
		`usage: migrate [-path=<path>] <command> [<args>]

the 'up' command builds the operations that will take the database from the state 
left in by the previous migration so that it is up-to-date with regard to this migration.

Commands:
   up             	Builds the operations that will migrate the database 'up'
   reset			Resets migration index to '0' and resets the system with just the root user
   help           	Show this help

'-path' defaults to current working directory.
`)

}

// RunMigrationUpToDate parses all 'up' files, sorts them and exec in order
// if version already exist it skipts it the newer version.
func RunMigrationUpToDate(pipe chan interface{}, mPath string) {
	d, files, version, err := initDriverAndReadMigrationFiles(mPath)
	if err != nil {
		go ClosePipeChannel(pipe, err)
		return
	}

	migrationFiles, err := files.SortByVersion(version)
	if err != nil {
		if err2 := d.Close(); err2 != nil {
			pipe <- err2
		}
		go ClosePipeChannel(pipe, err)
		return
	}

	// before migrating make sure the user did not interrupt!
	signals := handleInterrupts()
	defer signal.Stop(signals)

	if len(migrationFiles) > 0 {
		for _, f := range migrationFiles {
			pipe1 := NewPipeChannel()
			go d.Migrate(f, pipe1)
			if ok := WaitAndRedirect(pipe1, pipe, signals); !ok {
				break
			}
		}
		if err := d.Close(); err != nil {
			pipe <- err
		}
		go ClosePipeChannel(pipe, nil)
		return
	}
	if err := d.Close(); err != nil {
		pipe <- err
	}
	go ClosePipeChannel(pipe, nil)
	return
}

func initDriverAndReadMigrationFiles(mPath string) (Driver, *MFiles, uint64, error) {
	d, err := NewDriver()
	if err != nil {
		return nil, nil, 0, err
	}
	files, err := ReadFiles(mPath, FilenameRegex())
	if err != nil {
		d.Close()
		return nil, nil, 0, err
	}
	version, err := d.Version()
	if err != nil {
		d.Close()
		return nil, nil, 0, err
	}
	return d, &files, version, nil
}

//--------------------Start - of - PostgresDriver-----------------------
func (dbDriver *DBDriver) Initialize() error {
	dbDriver.Schema = Schema

	//TODO:: config should have a helper function to return connections
	cfx, err := config.GetDBConfig()
	if err != nil {
		return fmt.Errorf("error: failed to get db config %v", err)
	}

	_, connectionString := cfx.GetConnectionString()
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	dbDriver.db = db
	dbDriver.tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("unable to create sql transaction: %s", err.Error())
	}

	return initMigrationTables(dbDriver)
}

func initMigrationTables(dbDriver *DBDriver) error {
	if dbDriver.Schema != "" {
		if _, err := dbDriver.tx.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", dbDriver.Schema)); err != nil {
			return fmt.Errorf("initMigrationTables create schema err =%s", err)
		}

		if _, err := dbDriver.tx.Exec(fmt.Sprintf("SET search_path = %s, pg_catalog;", dbDriver.Schema)); err != nil {
			return fmt.Errorf("initMigrationTables set search_path err =%s", err)
		}
	}

	r := dbDriver.tx.QueryRow("SELECT count(*) FROM information_schema.tables WHERE table_name = $1 AND table_schema = (SELECT current_schema());", tableName)
	c := 0
	if err := r.Scan(&c); err != nil {
		return err
	}
	if c > 0 {
		return nil
	}
	if _, err := dbDriver.tx.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (version bigint not null primary key);"); err != nil {
		return err
	}
	return nil
}

//Close driver handler connections
func (dbDriver *DBDriver) Close() error {
	if err := dbDriver.tx.Commit(); err != nil {
		dbDriver.db.Close()
		return err
	}
	return dbDriver.db.Close()
}

func Reset() error {
	d, err := NewDriver()
	if err != nil {
		return err
	}

	if _, err = d.ResetVersion(); err != nil {
		_ = d.Close()
		return err
	}

	_ = d.Close()
	return nil
}

func (dbDriver *DBDriver) ResetVersion() (uint64, error) {
	var retVal uint64
	err := dbDriver.tx.QueryRow("DELETE FROM " + tableName).Scan(&retVal)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return retVal, nil
	}
}

func (dbDriver *DBDriver) Migrate(f File, pipe chan interface{}) {
	defer close(pipe)
	pipe <- f

	tx := dbDriver.tx
	if f.Direction == UpDirection {
		if _, err := tx.Exec("INSERT INTO "+tableName+" (version) VALUES ($1)", f.Version); err != nil {
			pipe <- err
			if err := tx.Rollback(); err != nil {
				pipe <- err
			}
			return
		}
	}

	if err := f.GetContent(); err != nil {
		pipe <- err
		return
	}

	if _, err := tx.Exec(string(f.Content)); err != nil {
		switch pqErr := err.(type) {
		case *pq.Error:
			offset, err := strconv.Atoi(pqErr.Position)
			if err == nil && offset >= 0 {
				lineNo, columnNo := LineColumnFromOffset(f.Content, offset-1)
				errorPart := LinesBeforeAndAfter(f.Content, lineNo, 5, 5, true)
				pipe <- fmt.Errorf("%s %v: %s in line %v, column %v:\n\n%s", pqErr.Severity, pqErr.Code, pqErr.Message, lineNo, columnNo, string(errorPart))
			} else {
				pipe <- fmt.Errorf("%s %v: %s", pqErr.Severity, pqErr.Code, pqErr.Message)
			}

			if err := tx.Rollback(); err != nil {
				pipe <- err
			}
			return
		default:
			pipe <- err
			if err := tx.Rollback(); err != nil {
				pipe <- err
			}
			return
		}
	}
}

//Version return migrate version of current schema
func (dbDriver *DBDriver) Version() (uint64, error) {
	var version uint64
	err := dbDriver.tx.QueryRow("SELECT version FROM " + tableName + " ORDER BY version DESC LIMIT 1").Scan(&version)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return version, nil
	}
}

// GetContent reads file content
func (f *File) GetContent() error {

	if len(f.Content) == 0 {
		content, err := ioutil.ReadFile(path.Join(f.Path, f.FileName))
		if err != nil {
			return err
		}
		f.Content = content
	}
	return nil
}

// SortByVersion sorts and validates migration files
func (mf *MFiles) SortByVersion(version uint64) (Files, error) {
	sort.Sort(mf)
	files := make(Files, 0)
	for _, migrationFile := range *mf {
		if migrationFile.Version > version && migrationFile.UpFile != nil {
			files = append(files, *migrationFile.UpFile)
		}
	}
	return files, nil
}

// ReadFiles from path directory. files must match the prefix require
func ReadFiles(path string, filenameRegex *regexp.Regexp) (files MFiles, err error) {
	// find all migration files in path
	ioFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	type tmpFile struct {
		version  uint64
		name     string
		filename string
		d        Direction
	}
	tmpFiles := make([]*tmpFile, 0)
	tmpFileMap := map[uint64]map[Direction]tmpFile{}
	for _, file := range ioFiles {
		version, name, d, err := filterFileName(file.Name(), filenameRegex)
		if err == nil {
			if _, ok := tmpFileMap[version]; !ok {
				tmpFileMap[version] = map[Direction]tmpFile{}
			}
			if existing, ok := tmpFileMap[version][d]; !ok {
				tmpFileMap[version][d] = tmpFile{version: version, name: name, filename: file.Name(), d: d}
			} else {
				return nil, fmt.Errorf("duplicate migration file version %d : %q and %q", version, existing.filename, file.Name())
			}
			tmpFiles = append(tmpFiles, &tmpFile{version, name, file.Name(), d})
		}
	}

	tVersions := make(map[uint64]bool)
	newFiles := make(MFiles, 0)
	for _, file := range tmpFiles {
		if _, ok := tVersions[file.version]; !ok {
			migrationFile := MFile{
				Version: file.version,
			}

			switch file.d {
			case UpDirection:
				migrationFile.UpFile = &File{
					Path:      path,
					FileName:  file.filename,
					Version:   file.version,
					Name:      file.name,
					Content:   nil,
					Direction: UpDirection,
				}

			default:
				return nil, errors.New("unsupported direction type")
			}

			newFiles = append(newFiles, migrationFile)
			tVersions[file.version] = true
		}
	}

	sort.Sort(&newFiles)
	return newFiles, nil
}

func filterFileName(filename string, filenameRegex *regexp.Regexp) (version uint64, name string, d Direction, err error) {
	matches := filenameRegex.FindStringSubmatch(filename)
	if len(matches) != 4 {
		return 0, "", 0, errors.New("Unable to parse filename schema")
	}

	version, err = strconv.ParseUint(matches[1], 10, 0)
	if err != nil {
		return 0, "", 0, fmt.Errorf("Unable to parse version '%v' in filename schema", matches[0])
	}

	if matches[3] == "up" {
		d = UpDirection
	} else if matches[3] == "down" {
		d = DownDirection
	} else {
		return 0, "", 0, fmt.Errorf("Unable to parse up|down '%v' in filename schema", matches[3])
	}

	return version, matches[2], d, nil
}

func LineColumnFromOffset(data []byte, offset int) (line, column int) {
	fs := token.NewFileSet()
	tf := fs.AddFile("", fs.Base(), len(data))
	tf.SetLinesForContent(data)
	pos := tf.Position(tf.Pos(offset))
	return pos.Line, pos.Column
}

func FilenameRegex() *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(fileExpression, fileExtension))
}

func LinesBeforeAndAfter(data []byte, line, before, after int, lineNumbers bool) []byte {
	startLine := line - before
	endLine := line + after
	lines := bytes.SplitN(data, []byte("\n"), endLine+1)

	if startLine < 0 {
		startLine = 0
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}

	selectLines := lines[startLine:endLine]
	newLines := make([][]byte, 0)
	lineCounter := startLine + 1
	lineNumberDigits := len(strconv.Itoa(len(selectLines)))
	for _, l := range selectLines {
		lineCounterStr := strconv.Itoa(lineCounter)
		if len(lineCounterStr)%lineNumberDigits != 0 {
			lineCounterStr = strings.Repeat(" ", lineNumberDigits-len(lineCounterStr)%lineNumberDigits) + lineCounterStr
		}

		lNew := l
		if lineNumbers {
			lNew = append([]byte(lineCounterStr+": "), lNew...)
		}
		newLines = append(newLines, lNew)
		lineCounter++
	}

	return bytes.Join(newLines, []byte("\n"))
}

// NewDriver creates a news instance db handler
func NewDriver() (Driver, error) {
	d := GetDriver(Scheme)
	if d == nil {
		return nil, fmt.Errorf("Driver '%s' not found", Scheme)
	}

	if fileExtension == "" {
		panic(fmt.Sprintf("%s returns empty string.", Scheme))
	}

	if fileExtension[0:1] == "." {
		panic(fmt.Sprintf("%s returned string must not start with a dot.", Scheme))
	}

	if err := d.Initialize(); err != nil {
		return nil, err
	}

	return d, nil
}

func RegisterDriver(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("driver: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetDriver(name string) Driver {
	driversMu.Lock()
	defer driversMu.Unlock()
	driver := drivers[name]
	return driver
}

// handleInterrupts ensures to monitor os singals and interrupts
func handleInterrupts() chan os.Signal {
	if interrupts {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		return c
	}
	return nil
}

// -----------------------Start -of - Pipe------------------------------------------------------
func NewPipeChannel() chan interface{} {
	return make(chan interface{}, 0)
}

func ClosePipeChannel(pipe chan interface{}, err error) {
	if err != nil {
		pipe <- err
	}
	close(pipe)
}

func WaitAndRedirect(pipe, redirectPipe chan interface{}, interrupt chan os.Signal) (ok bool) {
	errorReceived := false
	interruptsReceived := 0
	defer stopNotifyInterruptChannel(interrupt)
	if pipe != nil && redirectPipe != nil {
		for {
			select {

			case <-interrupt:
				interruptsReceived++
				if interruptsReceived > 1 {
					os.Exit(5)
				} else {
					// add white space at beginning for ^C splitting
					redirectPipe <- " Aborting after this migration ... Hit again to force quit."
				}

			case item, ok := <-pipe:
				if !ok {
					return !errorReceived && interruptsReceived == 0
				}
				redirectPipe <- item
				switch item.(type) {
				case error:
					errorReceived = true
				}
			}
		}
	}
	return !errorReceived && interruptsReceived == 0
}

func stopNotifyInterruptChannel(interruptChannel chan os.Signal) {
	if interruptChannel != nil {
		signal.Stop(interruptChannel)
	}
}

func writePipe(pipe chan interface{}) (ok bool) {
	okFlag := true
	if pipe != nil {
		for {
			select {
			case item, more := <-pipe:
				if !more {
					return okFlag
				}
				switch item.(type) {

				case string:
					fmt.Println(item.(string))

				case error:
					c := color.New(color.FgRed)
					c.Printf("%s\n\n", item.(error).Error())
					okFlag = false

				case File:
					f := item.(File)
					if f.Direction == UpDirection {
						c := color.New(color.FgGreen)
						c.Print(">")
					} else if f.Direction == DownDirection {
						c := color.New(color.FgRed)
						c.Print("<")
					}
					fmt.Printf(" %s\n", f.FileName)

				default:
					text := fmt.Sprint(item)
					fmt.Println(text)
				}
			}
		}
	}
	return okFlag
}
