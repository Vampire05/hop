package main

import (
	"encoding/json"
	"os"
)

type Request struct {
	Name    string            `json:"name"`
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
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

	original := requests[index]

	// Header kopieren
	clonedHeaders := make(map[string]string)
	for k, v := range original.Headers {
		clonedHeaders[k] = v
	}

	clone := Request{
		Name:    original.Name,
		URL:     original.URL,
		Method:  original.Method,
		Body:    original.Body,
		Headers: clonedHeaders,
	}

	if index == len(requests)-1 {
		requests = append(requests, clone)
	} else {
		requests = append(requests[:index+1], append([]Request{clone}, requests[index+1:]...)...)
	}

	saveRequests()
}

func moveRequestUp(index int) int {
	if index <= 0 || index >= len(requests) {
		return -1
	}

	// swap
	requests[index-1], requests[index] = requests[index], requests[index-1]

	saveRequests()
	return index - 1
}

func moveRequestDown(index int) int {
	if index < 0 || index >= len(requests)-1 {
		return -1
	}

	// swap
	requests[index], requests[index+1] = requests[index+1], requests[index]

	saveRequests()
	return index + 1
}
