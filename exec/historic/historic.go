package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goAnalyzeNGINX/shared"
	"os"
	"strings"
	"time"
)

var (
	inputFile  = "./access.log"     // nginx log file
	outputFile = "./countries.json" // output json
)

type dataStruct struct {
	Requests  map[string]int
	IPs       map[string]int
	Timestamp int64
}

func exportJSON(data dataStruct) {
	json, _ := json.Marshal(data)
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	file.Write(json)
}

func main() {
	var finalData dataStruct

	execTime := time.Now().UnixMilli()
	finalData.Timestamp = time.Now().Unix()

	fmt.Println("Script ran at", time.Now())

	// parse database
	database, err := shared.ParseDatabase()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// open log and create scanner
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error: ", err)
	}

	// parse log file and find country
	reqCount := 0

	requests_countries := make(map[string]int) // amount of requests per country
	ips_countries := make(map[string]int)      // amount of unique visitors per country

	// keep a map of already visited IPs for when we want to distinguish
	// between req from each country vs visitors from each country
	analyzed_ips := make(map[string]bool)

	for scanner.Scan() {
		ip := strings.Fields(scanner.Text())[0]
		reqCount++

		country := shared.FindCountry(ip, database)

		if country == "Unknown" {
			continue
		}

		requests_countries[country]++
		if !analyzed_ips[ip] {
			analyzed_ips[ip] = true
			ips_countries[country]++
		}
	}

	finalData.Requests = requests_countries
	finalData.IPs = ips_countries

	exportJSON(finalData)

	fmt.Println(reqCount, "requests processed in", (time.Now().UnixMilli() - execTime), "ms")
}
