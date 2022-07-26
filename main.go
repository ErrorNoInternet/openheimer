package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

var (
	version     string = "0.1.0"
	database    *diskv.Diskv
	ipChannel   chan string = make(chan string)
	scanQueue   chan string = make(chan string)
	pingWorkers int
	scanWorkers int
	pinging     bool
	scanning    bool
	pinged      int64
	scanned     int64
	valid       int64

	databasePath   string
	logFile        string
	ipFile         string
	startingIP     string
	timeout        int
	maxPingWorkers int
	maxScanWorkers int
	verbose        *bool
)

func main() {
	flag.StringVar(&databasePath, "database", "openheimer.db", "The database to store the results in")
	flag.StringVar(&logFile, "logFile", "openheimer.log", "The file to store the logs in")
	flag.StringVar(&ipFile, "ipFile", "", "The file to extract IP addresses from")
	flag.StringVar(&startingIP, "startingIP", "1.0.0.0", "The IP address to start scanning from")
	flag.IntVar(&timeout, "timeout", 5, "The amount of time to wait before timing out")
	flag.IntVar(&maxPingWorkers, "maxPingWorkers", 4000, "The maximum amount of workers to ping IPs")
	flag.IntVar(&maxScanWorkers, "maxScanWorkers", 1000, "The maximum amount of workers to scan IPs")
	verbose = flag.Bool("verbose", false, "Display everything that's happening")
	displayVersion := flag.Bool("version", false, "Display the current version of OpenHeimer")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("OpenHeimer v%v\n", version)
		return
	}

	startTime := time.Now().Unix()
	file, err := os.Create(logFile)
	if err != nil {
		log.Fatalf("Unable to create %v: %v\n", logFile, err.Error())
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	flatTransform := func(s string) []string { return []string{} }
	database = diskv.New(diskv.Options{
		BasePath:     databasePath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	go displayStatus()
	go pingIPs(ipChannel)
	go scanIPs(scanQueue)
	if ipFile != "" {
		result := readFromFile(ipFile, ipChannel)
		if result == 1 {
			return
		}
	} else {
		result := generateIPs(startingIP, ipChannel)
		if result == 1 {
			return
		}
	}

	for pinging || scanning {
		time.Sleep(1 * time.Second)
	}
	log.Printf("Done! Finished in %v seconds. Pinged: %v, Scanned: %v, Valid: %v\n", time.Now().Unix()-startTime, pinged, scanned, valid)
}

func displayStatus() {
	for {
		time.Sleep(5 * time.Second)
		log.Printf("Pinged: %v, Scanned: %v, Valid: %v\n", pinged, scanned, valid)
	}
}
