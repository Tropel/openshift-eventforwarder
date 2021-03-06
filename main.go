package main

import (
	"os"
	"net/http"
	"bufio"
	"log"
	"encoding/json"
	"golang.org/x/build/kubernetes/api"
	"fmt"
)

type Stream struct {
	Type string `json:"type,omitempty"`
	Event api.Event `json:"object"`
}

func main() {
	apiAddr := os.Getenv("OPENSHIFT_API_URL")
	apiToken := os.Getenv("OPENSHIFT_TOKEN")

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiAddr + "/api/v1/events?watch=true", nil)
	if (err != nil) {
		log.Fatal("Error while opening connection", err)
	}
	req.Header.Add("Authorization", "Bearer " + apiToken)
	resp, err := client.Do(req)

	if (err != nil) {
		log.Fatal("Error while connecting to:", apiAddr, err)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if (err != nil) {
			log.Fatal("Error reading from response stream.", err)
		}

		event := Stream{}
		decErr := json.Unmarshal(line, &event)
		if (decErr != nil) {
			log.Fatal("Error decoding json", err)
		}

		fmt.Printf("Project: %v, Time: %v | Name: %v | Kind: %v | Reason: %v | Message: %v",
			event.Event.Namespace, event.Event.LastTimestamp, event.Event.Name,
			event.Event.Kind, event.Event.Reason, event.Event.Message)
	}
}