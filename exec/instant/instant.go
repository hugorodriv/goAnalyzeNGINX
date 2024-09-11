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
	"sync/atomic"
)

var (
	reqLimit          = 100
	clientLimit int32 = 10
	inputFile         = "./data/fakeLogs.log" // nginx log file
)

func listen(database shared.IpDatabase, atomicClients *atomic.Int32, w http.ResponseWriter, f http.Flusher, ctx <-chan struct{}) {
	cmd := exec.Command("tail", "--lines", "1", "--retry", "--follow", inputFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)

	go func(atomicClients *atomic.Int32) {
		counter := 1
		for scanner.Scan() {
			if counter > reqLimit {
				fmt.Fprintf(w, "%s\n", "-")
				f.Flush()
				atomicClients.Add(-1)
				return
			}
			select {
			case <-ctx:
				atomicClients.Add(-1)
				fmt.Printf("Client disconnected\n")
				return
			default:
				text := scanner.Text()
				fields := strings.Fields(text)
				fmt.Fprintf(w, "%s", shared.FindCountry(fields[0], database))
				f.Flush()
			}
			counter++
		}
	}(atomicClients)

	err = cmd.Start()
	if err != nil {
		println(os.Stderr, "Error starting Cmd", err)
		atomicClients.Add(-1)
		fmt.Printf("Client disconnected\n")
		return
	}

	err = cmd.Wait()
	if err != nil {
		println(os.Stderr, "Error waiting for Cmd", err)
		atomicClients.Add(-1)
		fmt.Printf("Client disconnected\n")
		return
	}
}

func main() {
	var atomicClients atomic.Int32
	// Config SSE
	database, err := shared.ParseDatabase()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	http.HandleFunc("/instant_events", func(w http.ResponseWriter, r *http.Request) {
		if atomicClients.Load() > clientLimit {
			fmt.Println("Too many clients (>", clientLimit, ")")
			return
		}
		fmt.Printf("Clients: %d\n", atomicClients.Load())
		atomicClients.Add(1)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Encoding", "none")
		flusher, ok := w.(http.Flusher)
		if !ok {
			fmt.Println(ok)
		}

		ctx := r.Context().Done()
		listen(database, &atomicClients, w, flusher, ctx)
		<-r.Context().Done()
	})

	http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
