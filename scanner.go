package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PassTheMayo/mcstatus/v3"
)

func scanIps() {
	scanning = true
	var mutex sync.Mutex
	if *verbose {
		log.Println("Starting IP scanner...")
	}

	for {
		for scanWorkers >= maxScanWorkers {
			time.Sleep(100 * time.Millisecond)
		}

		address := <-scanQueue
		if address == "end" {
			break
		}
		ip := address
		var port uint16 = 25565
		if strings.Contains(address, ":") {
			segments := strings.Split(address, ":")
			ip = segments[0]
			parsedPort, _ := strconv.ParseUint(segments[1], 10, 16)
			port = uint16(parsedPort)
		}
		go scanIp(ip, port, &mutex)
		mutex.Lock()
		scanWorkers++
		scanned++
		mutex.Unlock()
	}

	for scanWorkers > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	scanning = false
}

func scanIp(ip string, port uint16, mutex *sync.Mutex) {
	if *verbose {
		log.Printf("Scanning %v:%v...\n", ip, port)
	}
	response, err := mcstatus.Status(ip, port, mcstatus.JavaStatusOptions{Timeout: time.Second * time.Duration(timeout)})
	if err != nil {
		if *verbose {
			log.Printf("Unable to scan %v:%v: %v\n", ip, port, err.Error())
		}
		mutex.Lock()
		scanWorkers--
		mutex.Unlock()
		return
	}
	jsonObject, err := json.Marshal(response)
	if err != nil {
		log.Printf("Unable to marshal server response: %v\n", err.Error())
		mutex.Lock()
		scanWorkers--
		mutex.Unlock()
		return
	}
	log.Printf("Found Minecraft server at %v:%v\n", ip, port)
	err = database.Write(fmt.Sprintf("%v:%v", ip, port), jsonObject)
	if err != nil {
		log.Printf("Unable to write to database: %v\n", err.Error())
	}

	mutex.Lock()
	scanWorkers--
	valid++
	mutex.Unlock()
}
