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
	maxGoroutines   int    = 10000
	activeScans     int
	lastScannedIp   string
	database        *bitcask.Bitcask
)

func log(text string, level int) {
	if level < minimumLogLevel {
		return
	}
	logLevel := "NONE"
	if level == -1 {
		logLevel = "OTHER"
	} else if level == 0 {
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
		log("Saving variables to database...", 0)
		errorObject := database.Put([]byte("last-ip"), []byte(lastScannedIp))
		database.Sync()
		if errorObject != nil {
			log("Unable to save variables: "+errorObject.Error(), 2)
		}
		log("Finished saving variables to database", 0)
	}
}

func main() {
	helpPage := fmt.Sprintf(
		"OpenHeimer v%v\n\n"+
			"\t-h, --help\t\tDisplay a list of available arguments\n"+
			"\t-s, --start\t\tRun OpenHeimer and start scanning IPs\n"+
			"\t-q, --query\t\tQuery something from the database\n"+
			"\t\tserver\t\tQuery a Minecraft server from the database\n"+
			"\t\tplayer\t\tQuery a Minecraft player from the database",
		version,
	)
	if len(os.Args) == 1 {
		fmt.Println(helpPage)
		return
	}
	showHelp := false
	startScanning := false
	queryData := false
	queryDataServer := false
	queryDataPlayer := false
	query := ""
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
			} else if argument == "player" && queryData == true {
				queryDataPlayer = true
			} else {
				if queryDataServer {
					query = argument
				} else if queryDataPlayer {
					query = argument
				} else {
					fmt.Println("Unknown argument: " + argument)
				}
			}
		}
	}
	if showHelp {
		fmt.Println(helpPage)
		return
	}
	if startScanning || queryData {
		log("Initializing database...", 1)
		success, errorObject := initializeDatabase()
		if !success {
			log("Unable to initialize database: "+errorObject.Error(), 2)
			return
		}
	}
	if startScanning {
		log("Launching autosave goroutine...", 1)
		go autosave()
		startOpenHeimer()
		return
	}
	if queryData == true {
		if queryDataServer == false && queryDataPlayer == false {
			fmt.Println("You need to specify something to query!")
			return
		} else if queryDataServer && query == "" {
			fmt.Println("List of found Minecraft servers:")
			for key := range database.Keys() {
				if string(key) != "last-ip" {
					fmt.Println(key)
				}
			}
		} else if queryDataServer && query != "" {
			serverData, errorObject := database.Get([]byte(query))
			if errorObject != nil {
				fmt.Println("Unable to query server: " + errorObject.Error())
				return
			}
			fmt.Println(serverData)
		}
	}
}

func startOpenHeimer() {
	log("Fetching last scanned IP from database...", 0)
	lastIpBytes, errorObject := database.Get([]byte("last-ip"))
	if errorObject != nil {
		log("Unable to fetch last scanned IP: "+errorObject.Error(), 0)
		lastIpBytes = []byte("1.0.0.0")
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

		for activeScans >= maxGoroutines {
			time.Sleep(500 * time.Millisecond)
		}
		activeScans += 1
		go sendPing(serverIp)

		segmentD += 1
		if segmentD > 255 {
			segmentD = 1
			segmentC += 1
			if segmentC > 255 {
				segmentC = 1
				segmentB += 1
				log("Scanning "+fmt.Sprintf("%v.%v.%v.%v", segmentA, segmentB, segmentC, segmentD)+"...", 0)
				if segmentB > 255 {
					segmentB = 1
					segmentA += 1
					if segmentA > 255 {
						segmentA = 0
						segmentB = 0
						segmentC = 0
						segmentD = 0
					}
				}
			}
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
	pinger := mcpinger.New(serverAddress, 25565, mcpinger.McPingerOption(mcpinger.WithTimeout(5*time.Second)))
	response, errorObject := pinger.Ping()
	activeScans -= 1
	if errorObject != nil {
		log(fmt.Sprintf("Unable to ping %v: %v", serverAddress, errorObject.Error()), -1)
		return
	}
	log(fmt.Sprintf(
		"%v running Minecraft %v (%v/%v): %v",
		serverAddress,
		response.Version.Name,
		response.Players.Online,
		response.Players.Max,
		response.Players.Sample,
	), 1)
	database.Put([]byte(serverAddress), []byte(fmt.Sprintf(
		"time:%v,version:%v,protocol:%v,motd:%v,players_online:%v,players_max:%v,sample:%v",
		time.Now().Unix(),
		response.Version.Name,
		response.Version.Protocol,
		response.Description.Text,
		response.Players.Online,
		response.Players.Max,
		response.Players.Sample,
	)))
	database.Sync()
}
