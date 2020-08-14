package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"

	"github.com/Go/azuremonitor/db/dbcontext"
	"github.com/fatih/color"
)

var version = flag.Bool("version", false, "Show version")
var filepath = flag.String("filepath", "", "")
var interrupts = true

const (
	Version = "1.0.0"
)

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

	case "all":
		pipe := NewPipeChannel()
		IngestApplication(pipe)
		os.Exit(1)

	case "azuremonitor":
		pipe := NewPipeChannel()
		go IngestApplication(pipe)
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
		`usage: ingest [-path=<path>] <command> [<args>]
Commands:
	all      		imports all .csv files
   	users         	imports users.csv into users
   	help           	Show this help

'-path' defaults to current working directory.
`)

}

func removeWhiteSpaces(input string) string {
	lead := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	mid := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	retVal := lead.ReplaceAllString(input, "")
	retVal = mid.ReplaceAllString(retVal, " ")
	return retVal
}

// IngestAzuremonitorAzuremonitor ingest azuremonitor' table
func IngestApplication(pipe chan interface{}) {
	csvFile, _ := os.Open("application.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	application := &dbcontext.Application{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			go ClosePipeChannel(pipe, err)
			log.Fatal(err)
			return
		}

		if line[0] == "applicationID" {
			log.Printf("..................")
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatalf("Error while casting application from csv file %v ", err)
		}
		application.Applicationid = id
		application.SubscriptionID = &line[1]
		application.Name = &line[2]
		application.TenantID = &line[3]
		application.GrantType = &line[4]
		application.ClientID = &line[5]
		application.ClientSecret = &line[6]

		err = application.Insert()
		if err != nil {
			log.Printf("Warning::: Insert Application %v", err)
		}

		log.Printf("successfully entered a new Application...%s [%s]", application.Applicationid, *application.Name)
	}

	// adding interrupts watcher
	signals := handleInterrupts()
	defer signal.Stop(signals)

	//go ClosePipeChannel(pipe, nil)
	return
}

// -----------------------Start -of - Pipe------------------------------------------------------

func handleInterrupts() chan os.Signal {
	if interrupts {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		return c
	}
	return nil
}

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

				default:
					text := fmt.Sprint(item)
					fmt.Println(text)
				}
			}
		}
	}
	return okFlag
}
