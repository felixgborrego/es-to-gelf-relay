package main

import (
	"bufio"
	"encoding/json"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"io"
	"log"
	"net/http"
	"os"
)

type ESLog struct {
	Log string
}

func main() {

	graylogAddr, ok := os.LookupEnv("GRAYLOG_ADDRESS")
	if ok {
		gelfWriter, err := gelf.NewTCPWriter(graylogAddr)
		if err != nil {
			log.Fatalf("gelf.NewWriter: %s", err)
		}

		// log to both stderr and graylog2
		log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
		log.Printf("Sending logs to '%s'", graylogAddr)
	} else {
		log.Println("GRAYLOG_ADDRESS env missing using stdout")
	}

	ch := make(chan string, 100)

	httphandler := func(w http.ResponseWriter, req *http.Request) {
		// ES messages may be send logs in batches
		reader := bufio.NewReader(req.Body)

		for {
			line, err := reader.ReadString('\n')
			if len(line) == 0 {
				break
			}

			// Process the json line
			var log ESLog
			json.Unmarshal([]byte(line), &log)
			if len(log.Log) > 0 {
				ch <- log.Log
			}

			if err != nil {
				break
			}
		}
	}

	go relay(ch)

	http.HandleFunc("/", httphandler)
	log.Println("Listing for ES messages at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func relay(ch <-chan string) {
	for {
		logLine := <-ch
		log.Println(logLine)
	}
}
