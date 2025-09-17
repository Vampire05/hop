package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var reset = "\033[0m"
var white = "\033[37m"
var invert = "\033[47m\033[30m%-*s\033[0m"

var pageDevider int = 5                       // Wo wird die Seite vertikal aufgeteilt
var contentStartRow int = 8                   // Wo startet die erste Zeile nach der Überschrift
var width int = 80                            // screen width
var height int = 25                           // screen height
var leftWidth int = (width / pageDevider) - 1 // Platz links
var menuHeight int = height - 10
var menuWidth int = width - 2
var oldState *term.State = nil
var startOnRight int = leftWidth + 4
var startOnLeft int = 3

func switchTerminalRaw() {
	// Terminal in Raw Mode schalten
	oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	fmt.Println("Old terminal state:", oldState)
}

func restoreTerminalMode() {
	fmt.Print(white)
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func readScreenDimensions() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	width = w
	height = h
	leftWidth = (w / pageDevider) - 1 // Platz links
	menuHeight = height - 10
	menuWidth = width - 2
	startOnRight = leftWidth + 4

	if err != nil {
		fmt.Println("Error reading screen dimensions:", err)
		os.Exit(0)
	}
}

func printBanner() {
	banner := `
    ╦ ╦╔═╗╔═╗       ╦ ╦┬ ┬┌─┐┌─┐┬─┐┌┬┐┌─┐─┐ ┬┌┬┐  ╔═╗┌─┐┌─┐┬─┐┌─┐┌┬┐┬┌─┐┌┐┌  ╔═╗┬  ┌─┐┬ ┬┌─┐┬─┐┌─┐┬ ┬┌┐┌┌┬┐
    ╠═╣║ ║╠═╝  ───  ╠═╣└┬┘├─┘├┤ ├┬┘ │ ├┤ ┌┴┬┘ │   ║ ║├─┘├┤ ├┬┘├─┤ │ ││ ││││  ╠═╝│  ├─┤└┬┘│ ┬├┬┘│ ││ ││││ ││
    ╩ ╩╚═╝╩         ╩ ╩ ┴ ┴  └─┘┴└─ ┴ └─┘┴ └─ ┴   ╚═╝┴  └─┘┴└─┴ ┴ ┴ ┴└─┘┘└┘  ╩  ┴─┘┴ ┴ ┴ └─┘┴└─└─┘└─┘┘└┘─┴┘
`
	bannerSmall := `
    ╦ ╦╔═╗╔═╗
    ╠═╣║ ║╠═╝
    ╩ ╩╚═╝╩  
`
	if width < 110 {
		fmt.Println(bannerSmall)
	} else {
		fmt.Println(banner)
	}

}

// Cursorposition setzen (1-basiert: row, col)
func locate(row, col int) {
	fmt.Printf("\x1b[%d;%dH", row, col)
}

// Bildschirm löschen, plattformunabhängig
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Fensterrahmen mit Überschrift, Linie darunter und 2-Spalten-Layout
func drawMain(withSeperator, withHints bool) {
	readScreenDimensions()
	clearScreen()

	fmt.Println(green)

	// obere Linie
	fmt.Print("╔")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("═")
	}
	fmt.Println("╗")

	// Überschrift zentrieren
	printBanner()

	// Linie unter Überschrift
	fmt.Print("╠")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("═")
	}
	fmt.Println("╣")

	// Restlicher Rahmen mit 2 Spalten
	for j := 0; j < menuHeight; j++ {
		fmt.Print("║")

		// linke Spalte
		for i := 0; i < leftWidth; i++ {
			fmt.Print(" ")
		}

		if withSeperator {
			// vertikale Trennung
			fmt.Print("│")
		} else {
			fmt.Print(" ")
		}

		// rechte Spalte
		for i := 0; i < (menuWidth - leftWidth - 1); i++ { // -2 Rahmen, -1 Trennlinie
			fmt.Print(" ")
		}

		fmt.Println("║") // rechte Rahmenlinie
	}

	// untere Linie
	fmt.Print("╚")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("═")
	}
	fmt.Println("╝")

	if withHints {
		locate(height+1, 1)
		fmt.Print(blue, "Mit ↑/↓ move, F1 = help")
	}

}

// writeListLeft zeigt links entweder Requests oder Editierfelder, abhängig davon, was übergeben wird
func writeListLeft(selected int, list []string) {
	fmt.Println(white)

	// linke Spalte: editierbare Felder
	for i, field := range list {
		locate(contentStartRow+i, startOnLeft)
		if i == selected {
			fmt.Printf(invert, leftWidth-1, field)
		} else {
			fmt.Printf("%-*s", leftWidth-1, field)
		}
	}

}

func writeAnyKeyHint() {
	locate(height, 1)
	fmt.Print(white, "--- Press any key to return ---")
	os.Stdin.Read(make([]byte, 1))
}

// writeContentRight zeigt rechts entweder den gesamten Request oder nur ein Feld
func writeContentRight(content Request, selected int, editMode bool, editFields []string) {
	fmt.Println(white)

	locate(contentStartRow, startOnRight)
	fmt.Print(yellow, "NAME: ", white, content.Name)

	locate(contentStartRow+2, startOnRight)
	fmt.Print(yellow, "METHOD: ", white, content.Method)

	locate(contentStartRow+3, startOnRight)
	fmt.Print(yellow, "URL: ", white, content.URL)

	locate(contentStartRow+5, startOnRight)
	fmt.Print(yellow, "BODY: ", white, content.Body)

	// HEADERS unterhalb vom Body
	row := contentStartRow + 7
	locate(row, startOnRight)
	fmt.Print(yellow, "HEADERS:")
	for k, v := range content.Headers {
		row++
		locate(row, startOnRight+2)
		fmt.Printf("%s: %s", white+k, v)
	}

	if editMode {
		field := editFields[selected]
		switch field {
		case "NAME":
			locate(contentStartRow, startOnRight+len("NAME")+2)
			fmt.Printf(invert, len(content.Name)+1, content.Name)
		case "METHOD":
			locate(contentStartRow+2, startOnRight+len("METHOD")+2)
			fmt.Printf(invert, len(content.Method)+1, content.Method)
		case "URL":
			locate(contentStartRow+3, startOnRight+len("URL")+2)
			fmt.Printf(invert, len(content.URL)+1, content.URL)
		case "BODY":
			locate(contentStartRow+5, startOnRight+len("BODY")+2)
			fmt.Printf(invert, len(content.Body)+1, content.Body)
		case "HEADERS":
			locate(row+1, startOnRight)
			fmt.Print(invert, leftWidth, "[EDIT HEADERS WITH ENTER]")
		}
		locate(height, 0)
		fmt.Print(yellow, field, white, ": [PRESS ENTER TO EDIT]", white)
	} else {
		locate(height, 0)
	}
}

func writeContentHelp() {

	locate(contentStartRow+1, startOnLeft)
	fmt.Print(white, "hop is a simple commandline http client written in go.")
	locate(contentStartRow+2, startOnLeft)
	fmt.Print(white, "All http request are saved to a .json file.")

	locate(contentStartRow+5, startOnRight)
	fmt.Print(yellow, "e = ", white, "Enter edit mode for the selected request")

	locate(contentStartRow+6, startOnRight)
	fmt.Print(yellow, "Del = ", white, "Delete the selected request")

	locate(contentStartRow+7, startOnRight)
	fmt.Print(yellow, "Enter = ", white, "Send the selected request")

	locate(contentStartRow+8, startOnRight)
	fmt.Print(yellow, "c = ", white, "Clone the selected request")

	locate(contentStartRow+9, startOnRight)
	fmt.Print(yellow, "Page up = ", white, "Move selected request up")

	locate(contentStartRow+10, startOnRight)
	fmt.Print(yellow, "Page down = ", white, "Move selected request up")
}

func editField(r *Request, field string) {
	locate(height, 0)
	for i := 1; i < width; i++ {
		fmt.Print(white, " ")
	}
	locate(height, 0)
	fmt.Print(yellow, field, ": ", white)

	// Aktuellen Wert ermitteln
	var current string
	switch field {
	case "NAME":
		current = r.Name
	case "URL":
		current = r.URL
	case "METHOD":
		current = r.Method
	case "BODY":
		current = r.Body
	case "HEADERS":
		// für HEADERS kein Single-Value → bleibt leer oder man könnte Key=Value ausgeben
	}

	// aktuellen Wert als Starttext anzeigen
	fmt.Print(current)
	input := []rune(current)

	for {
		buf := make([]byte, 3)
		n, _ := os.Stdin.Read(buf)
		if n > 0 {
			switch buf[0] {
			case 13, 10: // Enter
				fmt.Println()
				goto done
			case 127: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					fmt.Print("\b \b")
				}
			default:
				input = append(input, rune(buf[0]))
				fmt.Print(string(buf[0]))
			}
		}
	}
done:
	switch field {
	case "NAME":
		r.Name = string(input)
	case "URL":
		r.URL = string(input)
	case "METHOD":
		r.Method = string(input)
	case "BODY":
		r.Body = string(input)
	case "HEADERS":
		// Format: Key=Value
		parts := string(input)
		var key, value string
		if idx := strings.Index(parts, "="); idx != -1 {
			key = parts[:idx]
			value = parts[idx+1:]
			if r.Headers == nil {
				r.Headers = make(map[string]string)
			}
			r.Headers[key] = value
		}
	}
}
