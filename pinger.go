package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func pingIPs(ipChannel chan string) {
	pinging = true
	var mutex sync.Mutex
	if *verbose {
		log.Println("Starting IP pinger...")
	}

	for {
		for pingWorkers >= maxPingWorkers {
			time.Sleep(100 * time.Millisecond)
		}

		ip := <-ipChannel
		if ip == "end" {
			break
		}
		go pingIP(ip, &mutex)
		mutex.Lock()
		pingWorkers++
		pinged++
		mutex.Unlock()
	}

	for pingWorkers > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	scanQueue <- "end"
	pinging = false
}

func pingIP(ip string, mutex *sync.Mutex) {
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "25565"), time.Second*time.Duration(timeout))
	if err != nil {
		if *verbose {
			log.Printf("Unable to ping %v: %v\n", ip, err.Error())
		}
		mutex.Lock()
		pingWorkers--
		mutex.Unlock()
		return
	}
	if connection != nil {
		defer connection.Close()
		scanQueue <- ip
	}

	mutex.Lock()
	pingWorkers--
	mutex.Unlock()
}
