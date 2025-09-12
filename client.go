package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func processRequest(r Request) {
	method := strings.TrimSpace(strings.ToUpper(r.Method))
	switch method {
	case "GET":
		fire_request("GET", r.URL, r.Body)
	case "POST":
		fire_request("POST", r.URL, r.Body)
	case "PUT":
		fire_request("PUT", r.URL, r.Body)
	case "DELETE":
		fire_request("DELETE", r.URL, r.Body)
	case "PATCH":
		fire_request("PATCH", r.URL, r.Body)
	default:
		fmt.Println("UNKNOWN HTTP METHOD")
	}
}

func fire_request(method string, url string, data string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()

	showResponse(*resp)
}

func showResponse(resp http.Response) {
	if resp.StatusCode == 200 {
		fmt.Println(green, "Response status:", resp.StatusCode)
	} else {
		fmt.Println(red, "Response status:", resp.StatusCode)
	}

	// Alle Header ausgeben
	fmt.Println(yellow, "Response Headers:")
	for key, values := range resp.Header {
		for _, v := range values {
			fmt.Printf("    %s: %s\n", key, v)
		}
	}

	// Body komplett ausgeben
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(white, scanner.Text(), reset)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(red, "ERROR beim Lesen des Bodys:", err, reset)
		return
	}
}
