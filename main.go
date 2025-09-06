package main

import (
	"wumpiwolf.de/hop/keyboardmanager"
)

func main() {
	loadRequests()
	runMenu()
}

func runMenu() {
	switchTerminalRaw()
	defer restoreTerminalMode()

	selected := 0
	currentRequest := 0
	editMode := false
	editFields := []string{"NAME", "METHOD", "URL", "BODY"}

	for {
		clearScreen()
		drawMain(true, true)

		if editMode {
			//writeListLeft(selected, editFields)
			writeContentRight(requests[currentRequest], selected, true, editFields)
		} else {
			var requestNames []string
			for _, v := range requests {
				requestNames = append(requestNames, v.Name)
			}
			writeListLeft(selected, requestNames)
			writeContentRight(requests[selected], 0, false, nil)
		}

		key := keyboardmanager.ReadKey() // ← alles über keyboardmanager

		switch key {
		case keyboardmanager.KeyEscape:
			if editMode {
				editMode = false
				selected = currentRequest
				saveRequests()
			} else {
				clearScreen()
				return
			}
		case "e":
			if !editMode {
				editMode = true
				currentRequest = selected
				selected = 0
			}
		case "c":
			if !editMode {
				cloneRequest(selected)
			}
		case keyboardmanager.KeyEnter:
			if editMode {
				editField(&requests[currentRequest], editFields[selected])
			} else {
				clearScreen()
				fire_get(requests[selected].URL)
				writeAnyKeyHint()
			}
		case keyboardmanager.KeyArrowUp:
			if selected > 0 {
				selected--
			}
		case keyboardmanager.KeyArrowDown:
			if editMode && selected < len(editFields)-1 {
				selected++
			} else if !editMode && selected < len(requests)-1 {
				selected++
			}
		case keyboardmanager.KeyDelete:
			deleteRequest(selected)
			if selected >= len(requests) && len(requests) > 0 {
				selected = len(requests) - 1
			}
		case keyboardmanager.KeyF1:
			drawMain(false, false)
			writeContentHelp()
			writeAnyKeyHint()
		}
	}
}
