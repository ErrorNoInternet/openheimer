package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
	"github.com/prologic/bitcask"
)

var (
	version         string = "1.0.0"
	minimumLogLevel int    = 0
	lastScannedIp   string
	database        *bitcask.Bitcask
)

func log(text string, level int) {
	if level < minimumLogLevel {
		return
	}
	logLevel := "NONE"
	if level == 0 {
		logLevel = "DEBUG"
	} else if level == 1 {
		logLevel = "INFO"
	} else if level == 2 {
		logLevel = "ERROR"
	}
	fmt.Println(fmt.Sprintf("[%v %v] %v", time.Now().Format("2006-01-02 15:04:05"), logLevel, text))
}

func autosave() {
	for {
		time.Sleep(10 * time.Second)
		log("Saving variables...", 1)
		errorObject := database.Put([]byte("last-ip"), []byte(lastScannedIp))
		database.Sync()
		if errorObject != nil {
			log("Unable to save variables: "+errorObject.Error(), 2)
		}
	}
}

func main() {
	helpPage := fmt.Sprintf(
		"OpenHeimer v%v\n\n"+
			"\t-h, --help\t\tDisplay a list of available arguments\n"+
			"\t-s, --start\t\tRun OpenHeimer and start scanning IPs\n"+
			"\t-q, --query\t\tQuery something from the database\n"+
			"\t\tserver\t\tQuery a Minecraft server from the database",
		version,
	)
	if len(os.Args) == 0 {
		fmt.Println(helpPage)
		return
	}
	showHelp := false
	startScanning := true
	queryData := false
	queryDataServer := false
	for index, argument := range os.Args {
		if index > 0 {
			if argument == "--help" || argument == "-h" {
				showHelp = true
			} else if argument == "--start" || argument == "-s" {
				startScanning = true
			} else if argument == "--query" || argument == "-q" {
				queryData = true
			} else if argument == "server" && queryData == true {
				queryDataServer = true
			} else {
				fmt.Println("Unknown argument: " + argument)
			}
		}
	}
	if showHelp {
		fmt.Println(helpPage)
		return
	}
	if startScanning {
		log("Initializing database...", 1)
		success, errorObject := initializeDatabase()
		if !success {
			log("Unable to initialize database: "+errorObject.Error(), 2)
			return
		}
		log("Launching autosave goroutine...", 1)
		go autosave()
		startOpenHeimer()
		return
	}
	if queryData == true && queryDataServer == false {
		fmt.Println("You need to specify something to query!")
		return
	}
}

func startOpenHeimer() {
	log("Fetching last scanned IP from database...", 0)
	lastIpBytes, errorObject := database.Get([]byte("last-ip"))
	if errorObject != nil {
		log("Unable to fetch last scanned IP: "+errorObject.Error(), 0)
		lastIpBytes = []byte("0.0.0.0")
	}
	lastIp := string(lastIpBytes)
	log(fmt.Sprintf("Starting IP scan from %v...", lastIp), 1)
	segments := strings.Split(lastIp, ".")
	segmentA, _ := strconv.Atoi(segments[0])
	segmentB, _ := strconv.Atoi(segments[1])
	segmentC, _ := strconv.Atoi(segments[2])
	segmentD, _ := strconv.Atoi(segments[3])

	for {
		serverIp := fmt.Sprintf("%v.%v.%v.%v", segmentA, segmentB, segmentC, segmentD)
		lastScannedIp = serverIp

		segmentD += 1
		if segmentD > 255 {
			segmentD = 1
			segmentC += 1
		}
		if segmentC > 255 {
			segmentC = 1
			segmentB += 1
		}
		if segmentB > 255 {
			segmentB = 1
			segmentA += 1
		}
		if segmentA > 255 {
			segmentA = 0
			segmentB = 0
			segmentC = 0
			segmentD = 0
		}
	}
}

func initializeDatabase() (bool, error) {
	_, errorObject := os.Stat("database")
	if errorObject != nil {
		log("Unable to find existing database! Creating new one...", 2)
	}
	database, errorObject = bitcask.Open("database")
	if errorObject != nil {
		return false, errorObject
	}
	return true, nil
}

func sendPing(serverAddress string) {
	pinger := mcpinger.New(serverAddress, 25565)
	response, errorObject := pinger.Ping()
	if errorObject != nil {
		log(fmt.Sprintf("Unable to ping %v: %v", serverAddress, errorObject.Error()), 2)
		return
	}
	log(response.Version.Name, 1)
}
