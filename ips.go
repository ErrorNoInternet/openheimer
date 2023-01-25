package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func readFromFile(filePath string, outputChannel chan string) int {
	log.Printf("Reading IPs from %v...\n", filePath)
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read file: %v\n", err.Error())
		return 1
	}
	ips := strings.Split(strings.TrimSpace(string(fileData)), "\n")
	ok, index := validateIps(ips)
	if !ok {
		log.Fatalf("The IP address on line %v is invalid!\n", index+1)
		return 1
	}
	log.Printf("Successfully read %v IP addresses from file\n", len(ips))
	for _, ip := range ips {
		if strings.TrimSpace(ip) != "" {
			outputChannel <- ip
		}
	}

	time.Sleep(100 * time.Millisecond)
	outputChannel <- "end"
	return 0
}

func generateIps(startingIp string, outputChannel chan string) int {
	log.Printf("Generating IPs from %v...\n", startingIp)
	ok, _ := validateIps([]string{startingIp})
	if !ok {
		log.Println("The starting IP address you specified is invalid!")
		return 1
	}

	segments := strings.Split(startingIp, ".")
	segmentA, _ := strconv.Atoi(segments[0])
	segmentB, _ := strconv.Atoi(segments[1])
	segmentC, _ := strconv.Atoi(segments[2])
	segmentD, _ := strconv.Atoi(segments[3])
	for {
		serverIp := fmt.Sprintf("%v.%v.%v.%v", segmentA, segmentB, segmentC, segmentD)
		outputChannel <- serverIp

		segmentD += 1
		if segmentD > 255 {
			segmentD = 0
			segmentC += 1
			if segmentC > 255 {
				segmentC = 0
				segmentB += 1
				log.Printf("Scanning IP range %v.%v.*.*...\n", segmentA, segmentB)
				if segmentB > 255 {
					segmentB = 0
					segmentA += 1
					if segmentA > 255 {
						break
					}
				}
			}
		}
	}

	time.Sleep(100 * time.Millisecond)
	outputChannel <- "end"
	return 0
}

func validateIps(ips []string) (bool, int) {
	for index, ip := range ips {
		if strings.TrimSpace(ip) == "" {
			continue
		}
		if strings.Contains(ip, ":") {
			addressSegments := strings.Split(ip, ":")
			ip = addressSegments[0]
			port := addressSegments[1]

			segments := strings.Split(ip, ".")
			if len(segments) != 4 {
				return false, index
			}
			for _, segment := range segments {
				number, err := strconv.Atoi(segment)
				if err != nil {
					return false, index
				}
				if number < 0 || number > 255 {
					return false, index
				}
			}

			number, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				return false, index
			}
			if number < 0 || number > 65535 {
				return false, index
			}
		} else {
			segments := strings.Split(ip, ".")
			if len(segments) != 4 {
				return false, index
			}
			for _, segment := range segments {
				number, err := strconv.Atoi(segment)
				if err != nil {
					return false, index
				}
				if number < 0 || number > 255 {
					return false, index
				}
			}
		}
	}
	return true, -1
}
