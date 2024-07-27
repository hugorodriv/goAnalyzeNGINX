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
	inputFile  = "./data/access.log"     // nginx log file
	outputFile = "./data/countries.json" // output json
)

type dataStruct struct {
	Countries map[string]int
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
	_ = database
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
	var count int

	countries := make(map[string]int)
	for scanner.Scan() {
		ip := strings.Fields(scanner.Text())[0]
		count++

		country := shared.FindCountry(ip, database)

		if country == "Unknown" {
			continue
		}
		countries[country]++
	}

	finalData.Countries = countries

	exportJSON(finalData)

	fmt.Println(count, "requests processed in", (time.Now().UnixMilli() - execTime), "ms")
}
