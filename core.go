package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
)

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
	playerList := ""
	currentData, errorObject := database.Get([]byte(serverAddress))
	if errorObject == nil {
		segments := strings.Split(string(currentData), "\n")
		for _, segment := range segments {
			if strings.HasPrefix(segment, "players_sample:") {
				existingPlayers := strings.Split(strings.Split(segment, ":")[1], "|")
				for _, existingPlayer := range existingPlayers {
					exists := false
					for _, newPlayer := range players {
						if newPlayer == existingPlayer {
							exists = true
						}
					}
					if !exists {
						playerList += existingPlayer + "|"
					}
				}
			}
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
		strings.Trim(playerList+strings.Join(players, "|"), "|"),
	)))
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
			segmentD = 0
			segmentC += 1
			if segmentC > 255 {
				segmentC = 0
				segmentB += 1
				log("Scanning "+fmt.Sprintf("%v.%v.*.*", segmentA, segmentB)+"...", 1)
				if segmentB > 255 {
					segmentB = 0
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
