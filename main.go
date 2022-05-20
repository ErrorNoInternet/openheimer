package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
	"github.com/prologic/bitcask"
)

var (
	version         string = "0.0.1"
	minimumLogLevel int    = 0
	maxGoroutines   int
	maxTimeout      int
	startIpAddress  string
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
		time.Sleep(20 * time.Second)
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
			"\t\t--goroutines\t\tMaximum number of goroutines (10000)\n"+
			"\t\t--timeout\t\tMaximum number of seconds to wait (5 seconds)\n"+
			"\t\t--ip-address\t\tSpecify which IP address to start scanning from\n"+
			"\t-q, --query\t\tQuery something from the database\n"+
			"\t\t--server\t\tQuery a Minecraft server from the database\n"+
			"\t\t--player\t\tQuery a Minecraft player from the database",
		version,
	)
	if len(os.Args) == 1 {
		fmt.Println(helpPage)
		return
	}
	showHelp := false
	startScanning := false
	startGoroutineCount := false
	startTimeout := false
	startCustomIp := false
	queryData := false
	queryDataServer := false
	queryDataPlayer := false
	query := ""
	maxGoroutines = 10000
	maxTimeout = 5
	startIpAddress = ""
	for index, argument := range os.Args {
		if index > 0 {
			if argument == "--help" || argument == "-h" {
				showHelp = true
			} else if argument == "--start" || argument == "-s" {
				startScanning = true
			} else if argument == "--query" || argument == "-q" {
				queryData = true
			} else if argument == "--server" && queryData == true {
				queryDataServer = true
			} else if argument == "--player" && queryData == true {
				queryDataPlayer = true
			} else if argument == "--goroutines" && startScanning == true {
				startGoroutineCount = true
			} else if argument == "--timeout" && startScanning == true {
				startTimeout = true
			} else if argument == "--ip-address" && startScanning == true {
				startCustomIp = true
			} else {
				if queryDataServer {
					query = argument
				} else if queryDataPlayer {
					query = argument
				} else if startGoroutineCount {
					var errorObject error
					maxGoroutines, errorObject = strconv.Atoi(argument)
					if errorObject != nil {
						fmt.Println(fmt.Sprintf("Unable to parse \"%v\" as an integer: %v", maxGoroutines, errorObject.Error()))
						return
					}
				} else if startTimeout {
					var errorObject error
					maxTimeout, errorObject = strconv.Atoi(argument)
					if errorObject != nil {
						fmt.Println(fmt.Sprintf("Unable to parse \"%v\" as an integer: %v", maxGoroutines, errorObject.Error()))
						return
					}
				} else if startCustomIp {
					invalid := false
					if strings.Count(argument, ".") != 3 {
						invalid = true
					}
					segments := strings.Split(argument, ".")
					for _, segment := range segments {
						number, errorObject := strconv.Atoi(segment)
						if errorObject != nil {
							invalid = true
						} else {
							if number > 255 {
								invalid = true
							}
						}
					}
					if invalid {
						fmt.Println(fmt.Sprintf("The IP you specified (%v) is invalid!", argument))
						return
					}
					startIpAddress = argument
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
		log(fmt.Sprintf("Max Goroutines: %v, Connect Timeout: %v second(s)", maxGoroutines, maxTimeout), 1)
		log("Launching autosave goroutine...", 1)
		go autosave()
		startOpenHeimer()
		return
	}
	if queryData {
		if queryDataServer == false && queryDataPlayer == false {
			fmt.Println("You need to specify something to query!")
			return
		} else if queryDataServer && query == "" {
			fmt.Println("List of found Minecraft servers:")
			for key := range database.Keys() {
				if string(key) != "last-ip" {
					fmt.Println(string(key))
				}
			}
		} else if queryDataServer && query != "" {
			serverData, errorObject := database.Get([]byte(query))
			if errorObject != nil {
				fmt.Println("Unable to query server: " + errorObject.Error())
				return
			}
			fmt.Println(string(serverData))
		} else if queryDataPlayer && query == "" {
			fmt.Println("List of found Minecraft players:")
			playerList := []string{}
			for key := range database.Keys() {
				if string(key) != "last-ip" {
					serverData, _ := database.Get([]byte(key))
					segments := strings.Split(string(serverData), "\n")
					for _, segment := range segments {
						if strings.HasPrefix(segment, "players_sample:") {
							players := strings.Split(strings.Split(segment, ":")[1], "|")
							for _, player := range players {
								exists := false
								for _, existingPlayer := range playerList {
									if player == existingPlayer {
										exists = true
									}
								}
								if !exists {
									playerList = append(playerList, player)
								}
							}
						}
					}
				}
			}
			for _, player := range playerList {
				if player != "" {
					fmt.Println(player)
				}
			}
		} else if queryDataPlayer && query != "" {
			found := false
			for serverKey := range database.Keys() {
				if string(serverKey) != "last-ip" {
					serverData, _ := database.Get([]byte(serverKey))
					segments := strings.Split(string(serverData), "\n")
					for _, segment := range segments {
						if strings.HasPrefix(segment, "players_sample:") {
							players := strings.Split(strings.Split(segment, ":")[1], "|")
							for _, player := range players {
								if strings.HasPrefix(player, query+"-") {
									fmt.Println("Player: " + query)
									fmt.Println("Server IP: " + string(serverKey))
									found = true
									segments := strings.Split(string(serverData), "\n")
									for _, segment := range segments {
										key := strings.Split(segment, ":")[0]
										value := strings.Split(segment, ":")[1]
										timestamp, _ := strconv.Atoi(value)
										if key == "time" {
											fmt.Println("Time: " + time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05") + " (" + value + ")")
										}
										if key == "version" {
											fmt.Println("Server Version: " + value)
										}
										if key == "protocol" {
											fmt.Println("Server Protocol: " + value)
										}
										if key == "motd" {
											fmt.Println("MOTD: " + value)
										}
										if key == "players_online" {
											fmt.Println("Online Player Count: " + value)
										}
										if key == "players_max" {
											fmt.Println("Maximum Player Count: " + value)
										}
										if key == "players_sample" {
											fmt.Println("Players: " + strings.Replace(value, "|", ", ", -1))
										}
									}
								}
							}
						}
					}
				}
			}
			if !found {
				fmt.Println("Unable to query player: player not found")
			}
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
	if startIpAddress != "" {
		lastIp = startIpAddress
	}
	log(fmt.Sprintf("Starting IP scan from %v...", lastIp), 1)
	segments := strings.Split(lastIp, ".")
	segmentA, _ := strconv.Atoi(segments[0])
	segmentB, _ := strconv.Atoi(segments[1])
	segmentC, _ := strconv.Atoi(segments[2])
	segmentD, _ := strconv.Atoi(segments[3])

	for {
		serverIp := fmt.Sprintf("%v.%v.%v.%v", segmentA, segmentB, segmentC, segmentD)
		lastScannedIp = serverIp

		for runtime.NumGoroutine() >= maxGoroutines {
			time.Sleep(500 * time.Millisecond)
		}
		go sendPing(serverIp)

		segmentD += 1
		if segmentD > 255 {
			segmentD = 1
			segmentC += 1
			if segmentC > 255 {
				segmentC = 1
				segmentB += 1
				log("Scanning "+fmt.Sprintf("%v.%v.*.*", segmentA, segmentB)+"...", 1)
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
	pinger := mcpinger.New(serverAddress, 25565, mcpinger.McPingerOption(mcpinger.WithTimeout(time.Duration(maxTimeout)*time.Second)))
	response, errorObject := pinger.Ping()
	if errorObject != nil {
		log(fmt.Sprintf("Unable to ping %v: %v", serverAddress, errorObject.Error()), -1)
		return
	}
	players := []string{}
	for _, player := range response.Players.Sample {
		if strings.TrimSpace(player.Name) != "" && !strings.Contains(player.Name, "§") {
			players = append(players, player.Name+"-"+player.ID)
		}
	}
	log(fmt.Sprintf(
		"%v running Minecraft %v (%v/%v): %v",
		serverAddress,
		response.Version.Name,
		response.Players.Online,
		response.Players.Max,
		strings.Join(players, ", "),
	), 1)
	database.Put([]byte(serverAddress), []byte(fmt.Sprintf(
		"time:%v\nversion:%v\nprotocol:%v\nmotd:%v\nplayers_online:%v\nplayers_max:%v\nplayers_sample:%v",
		time.Now().Unix(),
		response.Version.Name,
		response.Version.Protocol,
		response.Description.Text,
		response.Players.Online,
		response.Players.Max,
		strings.Join(players, "|"),
	)))
}
