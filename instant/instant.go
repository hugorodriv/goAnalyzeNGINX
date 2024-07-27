// TODO:
//
//	auto locate a .csv file in data/ (also for historic)
//	separate shared functionality into a standalone shared package
package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

var (
	inputFile   = "./fakeLogs.log"                       // nginx log file
	outputFile  = "./data/countries.json"                // output json
	databaseLoc = "./data/dbip-country-lite-2024-04.csv" // CSV ordered IP database. Format: 'start_ip,end_ip,country_code'
)

func main() {
	cmd := exec.Command("tail", "--retry", "--follow", inputFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		println(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		println(os.Stderr, "Error waiting for Cmd", err)
		return
	}
}
