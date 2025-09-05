package main

import (
	"fmt"
	"os"
)

func main() {
	loadRequests()

	// Menü starten mit geladenen Requests
	runMenu(requests)

	// Debug-Ausgabe: was geladen wurde
	for _, v := range requests {
		fmt.Println(v.Name)
	}
}

// runMenu zeigt eine Liste links an und lässt mit ↑/↓ navigieren
func runMenu(requests []Request) {

	switchTerminalRaw()
	defer restoreTerminalMode()

	selected := 0       // Auswahl in der Liste (Request oder Feld)
	currentRequest := 0 // aktuell gewählter Request
	editMode := false   // false = normal, true = editiermodus
	editFields := []string{"Name", "URL", "Method", "Body"}

	for {
		clearScreen()
		drawMain()

		if editMode {
			// Editiermodus: linke Spalte = Felder des ausgewählten Requests
			writeListLeft(selected, nil, true, editFields)
			writeContentRight(requests[currentRequest], selected, true, editFields, leftWidth+5)
		} else {
			// Normalmodus: linke Spalte = Requests
			writeListLeft(selected, requests, false, nil)
			writeContentRight(requests[selected], 0, false, nil, leftWidth+5)
		}

		var buf [1]byte
		os.Stdin.Read(buf[:])

		switch buf[0] {
		case 'q':
			if editMode {
				editMode = false
				selected = currentRequest
			} else {
				clearScreen()
				return
			}
		case 'e':
			if !editMode {
				editMode = true
				currentRequest = selected
				selected = 0 // Start bei erstem Feld
			}
		case 13, 10: // Enter
			if editMode {
				editField(&requests[currentRequest], editFields[selected], leftWidth+5, contentStartRow+3)
			} else {
				clearScreen()
				fire_get(requests[selected].URL)
				fmt.Print(blue, "\n--- Press any key to return ---")
				os.Stdin.Read(make([]byte, 1))
			}
		case 27: // Pfeiltasten
			var seq [2]byte
			os.Stdin.Read(seq[:])
			if seq[0] == 91 {
				switch seq[1] {
				case 65: // ↑
					if selected > 0 {
						selected--
					}
				case 66: // ↓
					if editMode && selected < len(editFields)-1 {
						selected++
					} else if !editMode && selected < len(requests)-1 {
						selected++
					}
				}
			}
		}
	}
}
