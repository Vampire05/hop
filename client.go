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
	if method == "GET" {
		fire_get(r.URL)
	} else if method == "POST" {
		fire_post(r.URL, r.Body)
	} else if method == "PUT" {
		fire_put(r.URL, r.Body)
	} else if method == "DELETE" {
		fire_delete(r.URL, r.Body)
	} else if method == "PATCH" {
		fmt.Println("Not implemented yet")
	} else {
		fmt.Println("UNKNOWN HTTP METHOD")
	}
}

func fire_get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(red, "ERROR", err, reset)
		return
	}
	defer resp.Body.Close()

	showResponse(*resp)
}

func fire_post(url string, data string) {
	fmt.Println(blue, data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Header wie Insomnia
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "*/*")
	//req.Header.Set("User-Agent", "insomnia/2023.4.0")
	//req.Header.Set("Cookie", "session_id_constructor=...; session_id_timetrack=...")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()

	// Body komplett ausgeben
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}

func fire_put(url string, data string) {

}

func fire_delete(url string, data string) {

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
