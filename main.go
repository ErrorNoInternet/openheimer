package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prologic/bitcask"
)

var (
	version         string = "0.0.1"
	minimumLogLevel int    = 0
	maxGoroutines   int
	maxTimeout      int
	startIpAddress  string
	ipAddressFile   string
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
			"\t\t--file\t\tScan all the IP addresses specified in the file\n"+
			"\t\t--goroutines\t\tMaximum number of goroutines (10000)\n"+
			"\t\t--timeout\t\tMaximum number of seconds to wait (5 seconds)\n"+
			"\t\t--ip-address\t\tSpecify which IP address to start scanning from\n"+
			"\t-v, --verbose\t\tDisplay everthing that's currently happening\n"+
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
	startCustomFile := false
	queryData := false
	queryDataServer := false
	queryDataPlayer := false
	queriedDataServer := false
	queriedDataPlayer := false
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
			} else if argument == "--verbose" || argument == "-v" {
				minimumLogLevel = -1
			} else if argument == "--query" || argument == "-q" {
				queryData = true
			} else if argument == "--server" && queryData == true {
				queryDataServer = true
				queriedDataServer = true
			} else if argument == "--player" && queryData == true {
				queryDataPlayer = true
				queriedDataPlayer = true
			} else if argument == "--file" && startScanning == true {
				startCustomFile = true
			} else if argument == "--goroutines" && startScanning == true {
				startGoroutineCount = true
			} else if argument == "--timeout" && startScanning == true {
				startTimeout = true
			} else if argument == "--ip-address" && startScanning == true {
				startCustomIp = true
			} else {
				if queryDataServer {
					query = argument
					queryDataServer = false
				} else if queryDataPlayer {
					query = argument
					queryDataPlayer = false
				} else if startGoroutineCount {
					var errorObject error
					maxGoroutines, errorObject = strconv.Atoi(argument)
					if errorObject != nil {
						fmt.Println(fmt.Sprintf("Unable to parse \"%v\" as an integer: %v", maxGoroutines, errorObject.Error()))
						return
					}
					startGoroutineCount = false
				} else if startTimeout {
					var errorObject error
					maxTimeout, errorObject = strconv.Atoi(argument)
					if errorObject != nil {
						fmt.Println(fmt.Sprintf("Unable to parse \"%v\" as an integer: %v", maxGoroutines, errorObject.Error()))
						return
					}
					startTimeout = false
				} else if startCustomIp {
					if verifyIp(argument) {
						fmt.Println(fmt.Sprintf("The IP you specified (%v) is invalid!", argument))
						return
					}
					startIpAddress = argument
					startCustomIp = false
				} else if startCustomFile {
					_, errorObject := os.Stat(argument)
					if errorObject != nil {
						fmt.Println(fmt.Sprintf("Unable to access %v: %v", argument, errorObject.Error()))
						return
					}
					ipAddressFile = argument
					startCustomFile = false
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
		log(fmt.Sprintf("Max Goroutines: %v, Connection Timeout: %v second(s)", maxGoroutines, maxTimeout), 1)
		log("Launching autosave goroutine...", 1)
		go autosave()
		startOpenHeimer()
		return
	}
	if queryData {
		if queriedDataServer == false && queriedDataPlayer == false {
			fmt.Println("You need to specify something to query!")
			return
		} else if queriedDataServer && query == "" {
			fmt.Println("List of found Minecraft servers:")
			for key := range database.Keys() {
				if string(key) != "last-ip" {
					fmt.Println(string(key))
				}
			}
		} else if queriedDataServer && query != "" {
			serverData, errorObject := database.Get([]byte(query))
			if errorObject != nil {
				fmt.Println("Unable to query server: " + errorObject.Error())
				return
			}
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
		} else if queriedDataPlayer && query == "" {
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
		} else if queriedDataPlayer && query != "" {
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

func verifyIp(ip string) bool {
	if strings.Count(ip, ".") != 3 {
		return true
	}
	segments := strings.Split(ip, ".")
	for _, segment := range segments {
		number, errorObject := strconv.Atoi(segment)
		if errorObject != nil {
			return true
		} else {
			if number > 255 {
				return true
			}
		}
	}
	return false
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
