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

func deleteRequest(index int) {
	if index < 0 || index >= len(requests) {
		return // oder Fehler zurückgeben
	}
	requests = append(requests[:index], requests[index+1:]...)
	saveRequests()
}

func cloneRequest(index int) {
	if index < 0 || index >= len(requests) {
		return // Index ungültig
	}

	// Original-Request kopieren
	original := requests[index]
	clone := Request{
		Name:   original.Name,
		URL:    original.URL,
		Method: original.Method,
		Body:   original.Body,
	}

	// Slice erweitern und Clone direkt nach dem Original einfügen
	if index == len(requests)-1 {
		// Original war letztes Element → einfach anhängen
		requests = append(requests, clone)
	} else {
		// Insert in der Mitte
		requests = append(requests[:index+1], append([]Request{clone}, requests[index+1:]...)...)
	}

	saveRequests()
}
