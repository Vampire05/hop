package main

import (
	"encoding/json"
	"os"
)

type Request struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

var requests []Request
var fileName = "requests.json"

func loadRequests() {
	file, err := os.ReadFile(fileName)
	if err != nil {
		requests = []Request{}
		return
	}
	json.Unmarshal(file, &requests)
}

func saveRequests() {
	data, _ := json.MarshalIndent(requests, "", "  ")
	os.WriteFile(fileName, data, 0644)
}
