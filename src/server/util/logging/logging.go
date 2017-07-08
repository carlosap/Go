package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"github.com/util"
)

const (
	logLevelEnv = "LOG_LEVEL"
)

const (
	debugPrefix = "DEBUG - "
	infoPrefix  = "INFO - "
	warnPrefix  = "WARN - "
	errorPrefix = "ERROR - "
)

var (
	//LogLevelConfig is the level set by the env
	logLevelConfig = util.GetSetEnv(logLevelEnv, "WARN")
)

var (
	//CurrentLogLevel of the system
	CurrentLogLevel LogLevel
	out             = New(os.Stderr, os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func init() {
	switch strings.ToUpper(logLevelConfig) {
	case "ERROR":
		CurrentLogLevel = LError
	case "WARN":
		CurrentLogLevel = LWarn
	case "INFO":
		CurrentLogLevel = LInfo
	case "DEBUG":
		CurrentLogLevel = LDebug
	default:
		CurrentLogLevel = LError
	}
}

//LogLevel is a logging level
type LogLevel uint8

//Logging levels
const (
	LMandatory = LogLevel(1 << iota)
	LError
	LWarn
	LInfo
	LDebug
)

//LevelLogger logs according to indicated log level
type LevelLogger struct {
	errLog *log.Logger
	stdLog *log.Logger
}

//New gets a new LevelLogger
func New(errOut, stdOut io.Writer, prefix string, flag int) *LevelLogger {
	return &LevelLogger{
		errLog: log.New(errOut, prefix, flag),
		stdLog: log.New(stdOut, prefix, flag),
	}
}

//Log sends the format and the params to the underlying logger
func (l *LevelLogger) Log(level LogLevel, formattedString string, params ...interface{}) {
	if level <= CurrentLogLevel {
		thisLog := l.errLog
		//Send mandatory (access) logs to the stdout writer
		//otherwise use stderr writer
		switch level {
		case LMandatory:
			thisLog = l.stdLog
		default:
		}
		thisLog.Output(3, fmt.Sprintf(formattedString, params...))
	}
}

//Fatalf is equivalent to calling Errorf followed by os.Exit(1)
func Fatalf(formattedString string, params ...interface{}) {
	out.Log(LError, formattedString, params)
	os.Exit(1)
}

//Panic is equivalent to calling Errorf followed by panic(params)
func Panic(formattedString string, params ...interface{}) {
	s := fmt.Sprintf(formattedString, params...)
	out.Log(LError, formattedString, params)
	panic(s)
}

//Error log
func Error(err error) {
	out.Log(LError, err.Error())
}

//Mandatory always logs regardless of logging level
func Mandatory(formattedString string, params ...interface{}) {
	out.Log(LMandatory, formattedString, params...)
}

//Errorf log
func Errorf(formattedString string, params ...interface{}) {
	out.Log(LError, errorPrefix+formattedString, params...)
}

//Warn log
func Warn(formattedString string, params ...interface{}) {
	out.Log(LWarn, warnPrefix+formattedString, params...)
}

//Info log
func Info(formattedString string, params ...interface{}) {
	out.Log(LInfo, infoPrefix+formattedString, params...)
}

//Debug log
func Debug(formattedString string, params ...interface{}) {
	out.Log(LDebug, debugPrefix+formattedString, params...)
}
