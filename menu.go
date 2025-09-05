package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var reset = "\033[0m"
var white = "\033[37m"

var pageDevider int = 4                       // Wo wird die Seite vertikal aufgeteilt
var contentStartRow int = 8                   // Wo startet die erste Zeile nach der Überschrift
var width int = 80                            // screen width
var height int = 25                           // screen height
var leftWidth int = (width / pageDevider) - 1 // Platz links
var menuHeight int = height - 10
var menuWidth int = width - 2
var oldState *term.State = nil

func switchTerminalRaw() {
	// Terminal in Raw Mode schalten
	oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	fmt.Println("Old terminal state:", oldState)
}

func restoreTerminalMode() {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func readScreenDimensions() {

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	width = w
	height = h
	leftWidth = (w / pageDevider) - 1 // Platz links
	menuHeight = height - 10
	menuWidth = width - 2

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
	fmt.Println(banner)
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

// Zeichnet eine durchgehende Linie im Contentbereich
func drawContentLine() {
	fmt.Printf(green)

	for i := 0; i < menuWidth-leftWidth-5; i++ {
		fmt.Print("-")
	}
}

// Fensterrahmen mit Überschrift, Linie darunter und 2-Spalten-Layout
func drawMain() {
	readScreenDimensions()

	fmt.Println(green)

	// obere Linie
	fmt.Print("+")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Überschrift zentrieren
	printBanner()

	// Linie unter Überschrift
	fmt.Print("+")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Restlicher Rahmen mit 2 Spalten
	for j := 0; j < menuHeight; j++ {
		fmt.Print("|")

		// linke Spalte
		for i := 0; i < leftWidth; i++ {
			fmt.Print(" ")
		}

		// vertikale Trennung
		fmt.Print("|")

		// rechte Spalte
		for i := 0; i < (menuWidth - leftWidth - 1); i++ { // -2 Rahmen, -1 Trennlinie
			fmt.Print(" ")
		}

		fmt.Println("|") // rechte Rahmenlinie
	}

	// untere Linie
	fmt.Print("+")
	for i := 0; i < menuWidth; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	locate(height+1, 1)
	fmt.Print(blue, "Mit ↑/↓ move, q = quit, e = edit, n = new request, c = clone request, Enter = send request")
}

// writeListLeft zeigt links entweder Requests oder Editierfelder, abhängig davon, was übergeben wird
func writeListLeft(selected int, requests []Request, editMode bool, editFields []string) {
	fmt.Println(white)

	if editMode {
		// linke Spalte: editierbare Felder
		for i, field := range editFields {
			locate(contentStartRow+i, 2)
			if i == selected {
				fmt.Printf("\033[47m\033[30m%-*s\033[0m", leftWidth-1, field)
			} else {
				fmt.Printf("%-*s", leftWidth-1, field)
			}
		}
	} else {
		// linke Spalte: Requests
		for i, item := range requests {
			locate(contentStartRow+i, 2)
			text := item.Name
			if len(text) > leftWidth-2 {
				text = text[:leftWidth-2]
			}
			if i == selected {
				fmt.Printf(" \033[47m\033[30m%-*s\033[0m", leftWidth-1, text)
			} else {
				fmt.Printf(" %-*s", leftWidth-1, text)
			}
		}
	}
}

// writeContentRight zeigt rechts entweder den gesamten Request oder nur ein Feld
func writeContentRight(content Request, selected int, editMode bool, editFields []string, x int) {
	fmt.Println(white)

	if editMode {
		// nur das aktuell ausgewählte Feld anzeigen
		field := editFields[selected]
		locate(contentStartRow, x)
		fmt.Print(yellow, "Aktueller Wert: ", white)
		switch field {
		case "Name":
			fmt.Print(content.Name)
		case "URL":
			fmt.Print(content.URL)
		case "Method":
			fmt.Print(content.Method)
		case "Body":
			fmt.Print(content.Body)
		}
	} else {
		// kompletter Request
		locate(contentStartRow, x)
		fmt.Print(yellow, "NAME: ", white, content.Name)

		locate(contentStartRow+2, x)
		fmt.Print(yellow, "METHOD: ", white, content.Method)

		locate(contentStartRow+3, x)
		fmt.Print(yellow, "URL: ", white, content.URL)

		locate(contentStartRow+5, x)
		fmt.Print(yellow, "BODY: ", white, content.Body)
	}
}

// PopupMenu zeigt ein einfaches ASCII-Popup mit Optionen an und gibt die gewählte Option zurück
func PopupMenu(title string, options []string) int {
	readScreenDimensions()

	// Popup-Größe
	pWidth := 50
	pHeight := len(options) + 4
	startRow := (height - pHeight) / 2
	startCol := (width - pWidth) / 2

	selected := 0

	drawPopup := func() {
		// Rahmen
		locate(startRow, startCol)
		fmt.Print("+")
		for i := 0; i < pWidth-2; i++ {
			fmt.Print("-")
		}
		fmt.Println("+")

		// Titel
		locate(startRow+1, startCol+2)
		fmt.Print(title)

		// Optionen
		for i, opt := range options {
			locate(startRow+2+i, startCol+2)
			if i == selected {
				fmt.Printf("\033[47m\033[30m%-*s\033[0m", pWidth-4, opt) // Highlight
			} else {
				fmt.Printf("%-*s", pWidth-4, opt)
			}
		}

		// untere Linie
		locate(startRow+pHeight-1, startCol)
		fmt.Print("+")
		for i := 0; i < pWidth-2; i++ {
			fmt.Print("-")
		}
		fmt.Println("+")
	}

	switchTerminalRaw()
	defer restoreTerminalMode()

	for {
		drawPopup()

		var buf [3]byte
		n, _ := os.Stdin.Read(buf[:])
		if n == 1 {
			switch buf[0] {
			case 27: // Escape-Sequenzen für Pfeile
				var arrow [2]byte
				os.Stdin.Read(arrow[:])
				if arrow[0] == 91 {
					switch arrow[1] {
					case 65: // Up
						if selected > 0 {
							selected--
						}
					case 66: // Down
						if selected < len(options)-1 {
							selected++
						}
					}
				}
			case 13: // Enter
				return selected
			case 'q':
				return -1
			}
		}
	}
}

func editField(r *Request, field string, x, y int) {
	locate(y, x)
	fmt.Print("Neuer Wert: ")

	input := []rune{}
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
	case "Name":
		r.Name = string(input)
	case "URL":
		r.URL = string(input)
	case "Method":
		r.Method = string(input)
	case "Body":
		r.Body = string(input)
	}
}
