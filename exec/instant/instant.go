// TODO:
//
//	auto locate a .csv file in data/ (also for historic)
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
)

var (
	reqLimit  = 100
	inputFile = "./data/fakeLogs.log" // nginx log file
)

func listen(database shared.IpDatabase, w http.ResponseWriter, f http.Flusher, ctx <-chan struct{}) {
	// Parse File
	cmd := exec.Command("tail", "--lines", "1", "--retry", "--follow", inputFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)

	counter := 1
	go func() {
		for scanner.Scan() {
			if counter > reqLimit {
				return
			}
			select {
			case <-ctx:
				return
			default:
				text := scanner.Text()
				fields := strings.Fields(text)
				fmt.Fprintf(w, "data: %s\n\n", shared.FindCountry(fields[0], database))
				f.Flush()
			}
			counter++
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

func main() {
	// Config SSE
	database, err := shared.ParseDatabase()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			fmt.Println(ok)
		}

		ctx := r.Context().Done()
		listen(database, w, flusher, ctx)
		<-r.Context().Done()
	})

	http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
