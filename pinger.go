package main

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func pingIps() {
	pinging = true
	var mutex sync.Mutex
	if *verbose {
		log.Println("Starting IP pinger...")
	}

	for {
		for pingWorkers >= maxPingWorkers {
			time.Sleep(100 * time.Millisecond)
		}

		address := <-addressChannel
		if address == "end" {
			break
		}
		ip, port := address, "25565"
		if strings.Contains(address, ":") {
			segments := strings.Split(address, ":")
			ip = segments[0]
			port = segments[1]
		}
		go pingIp(ip, port, &mutex)
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

func pingIp(ip string, port string, mutex *sync.Mutex) {
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), time.Second*time.Duration(timeout))
	if err != nil {
		if *verbose {
			log.Printf("Unable to ping %v:%v: %v\n", ip, port, err.Error())
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
