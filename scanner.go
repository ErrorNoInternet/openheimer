package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/PassTheMayo/mcstatus/v3"
)

func scanIPs(scanQueue chan string) {
	scanning = true
	var mutex sync.Mutex
	if *verbose {
		log.Println("Starting IP scanner...")
	}

	for {
		for scanWorkers >= maxScanWorkers {
			time.Sleep(100 * time.Millisecond)
		}

		ip := <-scanQueue
		if ip == "end" {
			break
		}
		go scanIP(ip, &mutex)
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

func scanIP(ip string, mutex *sync.Mutex) {
	if *verbose {
		log.Printf("Scanning %v...\n", ip)
	}
	response, err := mcstatus.Status(ip, 25565, mcstatus.JavaStatusOptions{Timeout: time.Second * time.Duration(timeout)})
	if err != nil {
		if *verbose {
			log.Printf("Unable to scan %v: %v\n", ip, err.Error())
		}
		mutex.Lock()
		scanWorkers--
		mutex.Unlock()
		return
	}
	jsonObject, err := json.Marshal(response)
	if err != nil {
		log.Printf("Unable to marshal server response: %v\n", err.Error())
		return
	}
	log.Printf("Found Minecraft server at %v\n", ip)
	err = database.Write(ip, jsonObject)
	if err != nil {
		log.Printf("Unable to save server details: %v\n", err.Error())
		return
	}

	mutex.Lock()
	scanWorkers--
	valid++
	mutex.Unlock()
}
