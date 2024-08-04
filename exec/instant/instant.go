// TODO:
//
//	auto locate a .csv file in data/ (also for historic)
//	separate shared functionality into a standalone shared package
package main

import (
	"bufio"
	"fmt"
	"goAnalyzeNGINX/shared"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var inputFile = "./data/fakeLogs.log" // nginx log file

func listen(database shared.IpDatabase) {
	cmd := exec.Command("tail", "--lines", "1", "--retry", "--follow", inputFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			ip := strings.Fields(scanner.Text())[0]
			processIP(ip, database)
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

func processIP(ip string, database shared.IpDatabase) {
	country := shared.FindCountry(ip, database)
	fmt.Println(country)
	http.HandleFunc("/events", eventsHandler)
	http.ListenAndServe(":8080", nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("Event %d", i))
		time.Sleep(2 * time.Second)
		w.(http.Flusher).Flush()
	}

	// Simulate closing the connection
	closeNotify := w.(http.CloseNotifier).CloseNotify()
	<-closeNotify
}

func main() {
	// println(shared.IpToInt("127.0.0.1"))
	database, err := shared.ParseDatabase()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	listen(database)
}
