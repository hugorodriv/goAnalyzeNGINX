package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	inputFile   = "./access.log"                  // nginx log file
	outputFile  = "countries.json"                // output json
	databaseLoc = "dbip-country-lite-2024-04.csv" // CSV ordered IP database. Format: 'start_ip,end_ip,country_code'
)

type ipDatabase struct {
	arr   []ipRange
	count int
}
type ipRange struct {
	begin   uint32
	end     uint32
	country string
}

type dataStruct struct {
	Timestamp int64
	Countries map[string]int
}

func ipToInt(ip string) uint32 {
	split := strings.Split(ip, ".")
	if len(split) != 4 {
		fmt.Println("Error: Invalid IP")
		fmt.Println("IP: ", ip)
		return 0
	}
	part1, _ := strconv.Atoi(split[0])
	part2, _ := strconv.Atoi(split[1])
	part3, _ := strconv.Atoi(split[2])
	part4, _ := strconv.Atoi(split[3])

	return uint32((part1 << 24) + (part2 << 16) + (part3 << 8) + part4)
}

func findCountry(ip string, database ipDatabase) string {
	ipInt := ipToInt(ip)

	// binary search
	var l int = 0
	var r int = database.count - 1
	for l <= r {
		m := l + (r-l)/2
		if database.arr[m].begin <= ipInt && database.arr[m].end >= ipInt {
			return database.arr[m].country
		} else if database.arr[m].begin > ipInt {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return "Unknown"
}

func parseDatabase() (ipDatabase, error) {
	file, err := os.Open(databaseLoc)
	if err != nil {
		fmt.Println("Error: ", err)
		return ipDatabase{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error: ", err)
		return ipDatabase{}, err
	}

	var tempArr []ipRange
	var count int

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if strings.Contains(line[0], ":") {
			break
		}
		begin := ipToInt(line[0])
		end := ipToInt(line[1])
		country := line[2]

		tempIpR := ipRange{begin, end, country}
		tempArr = append(tempArr, tempIpR)
		count++
	}
	return ipDatabase{tempArr, count}, nil
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

	finalData.Timestamp = time.Now().Unix()

	fmt.Println("Script ran at", time.Now())

	// parse database
	database, err := parseDatabase()
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

		country := findCountry(ip, database)

		if country == "Unknown" {
			continue
		}
		countries[country]++
	}

	finalData.Countries = countries

	exportJSON(finalData)

	fmt.Println("Req. processed: ", count)
}
