package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var reset = "\033[0m"
var white = "\033[37m"

type Request struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

var (
	fileName       = "requests.json"
	requests       []Request
	selected       int  // Auswahl in der Liste
	detailSelected int  // Auswahl in der Detail-View: 0=Name,1=Method,2=URL,3=Body
	inEditPopup    bool // true, wenn Popup für Feld-Edit offen
)

func loadRequests() {
	data, err := os.ReadFile(fileName)
	if err != nil {
		requests = []Request{}
		return
	}
	json.Unmarshal(data, &requests)
}

func saveRequests() {
	data, _ := json.MarshalIndent(requests, "", "  ")
	_ = os.WriteFile(fileName, data, 0644)
}

// ---------- GUI ----------

var guiInitialized bool = false

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Header
	if v, err := g.SetView("header", 0, 0, maxX-1, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		printHeader(v)
	}

	// Detail-View
	if v, err := g.SetView("details", maxX/4+1, 5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//v.Title = " [Details] "
		v.Wrap = true
		printDetails(g, v)
	}

	// Liste
	if v, err := g.SetView("list", 0, 5, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " [" + fileName + "] "
		printList(v)
	}

	// Fokus beim ersten Layout
	if !guiInitialized {
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}
		guiInitialized = true
	}

	return nil
}

func printHeader(v *gocui.View) {
	v.Clear()
	fmt.Fprintln(v)
	banner1 := "╦ ╦╔═╗╔═╗       ╦ ╦┬ ┬┌─┐┌─┐┬─┐┌┬┐┌─┐─┐ ┬┌┬┐  ╔═╗┌─┐┌─┐┬─┐┌─┐┌┬┐┬┌─┐┌┐┌  ╔═╗┬  ┌─┐┬ ┬┌─┐┬─┐┌─┐┬ ┬┌┐┌┌┬┐"
	banner2 := "╠═╣║ ║╠═╝  ───  ╠═╣└┬┘├─┘├┤ ├┬┘ │ ├┤ ┌┴┬┘ │   ║ ║├─┘├┤ ├┬┘├─┤ │ ││ ││││  ╠═╝│  ├─┤└┬┘│ ┬├┬┘│ ││ ││││ ││"
	banner3 := "╩ ╩╚═╝╩         ╩ ╩ ┴ ┴  └─┘┴└─ ┴ └─┘┴ └─ ┴   ╚═╝┴  └─┘┴└─┴ ┴ ┴ ┴└─┘┘└┘  ╩  ┴─┘┴ ┴ ┴ └─┘┴└─└─┘└─┘┘└┘─┴┘"

	// ANSI-Farbcodes: Rot = \033[31m, Gelb = \033[33m, Reset = \033[0m
	fmt.Fprintln(v, red+banner1)    // rot
	fmt.Fprintln(v, yellow+banner2) // gelb
	fmt.Fprintln(v, white+banner3)  // normal
}

func printList(v *gocui.View) {
	v.Clear()
	fmt.Fprint(v, "\n\n")
	for i, r := range requests {
		if i == selected {
			// invertiert darstellen
			fmt.Fprintf(v, "\033[30;43m%s\033[0m\n", r.Name)
		} else {
			fmt.Fprintf(v, "%s \n", r.Name)
		}
	}
}

func printDetails(g *gocui.Gui, v *gocui.View) {
	v.Clear()
	//fmt.Fprint(v, "\n\n")
	if len(requests) == 0 || selected < 0 || selected >= len(requests) {
		fmt.Fprintln(v, "Keine Requests")
		return
	}
	r := requests[selected]

	fields := []struct {
		label string
		value string
	}{
		{"Name", r.Name},
		{"Method", r.Method},
		{"URL", r.URL},
		{"Body", r.Body},
	}

	for i, f := range fields {
		cv := g.CurrentView()
		if i == detailSelected && cv != nil && cv.Name() == "details" && !inEditPopup {
			fmt.Fprintf(v, "\033[30;43m%s: %s\033[0m\n\n", f.label, f.value) // Schwarz auf Gelb
		} else {
			fmt.Fprintf(v, "%s: %s\n\n", yellow+f.label, white+f.value)
		}

	}
}

// ---------- Actions ----------

func editRequest(g *gocui.Gui, v *gocui.View) error {
	if len(requests) == 0 || inEditPopup {
		return nil
	}

	detailSelected = 0 // erstes Feld auswählen

	// Fokus auf Details setzen
	if _, err := g.SetCurrentView("details"); err != nil {
		return err
	}

	if dv, err := g.View("details"); err == nil {
		printDetails(g, dv)
	}

	return nil
}

func exitEditRequest(g *gocui.Gui, v *gocui.View) error {
	// Auswahl im Detail zurücksetzen
	detailSelected = -1

	// Fokus zurück auf Liste
	if lv, err := g.View("list"); err == nil {
		g.SetCurrentView("list")
		printList(lv)
	}

	// Details neu zeichnen, ohne Invertierung
	if dv, err := g.View("details"); err == nil {
		printDetails(g, dv)
	}

	return nil
}

func deleteRequest(g *gocui.Gui, v *gocui.View) error {
	if len(requests) == 0 {
		return nil
	}

	// Aktuellen Request entfernen
	requests = append(requests[:selected], requests[selected+1:]...)

	// Auswahl anpassen
	if selected >= len(requests) {
		selected = len(requests) - 1
	}

	// Speichern
	saveRequests()

	// Views aktualisieren
	if lv, err := g.View("list"); err == nil {
		printList(lv)
	}
	if dv, err := g.View("details"); err == nil {
		printDetails(g, dv)
	}

	return nil
}

func moveRequestUp(g *gocui.Gui, v *gocui.View) error {
	if selected > 0 {
		// Swap mit vorherigem
		requests[selected], requests[selected-1] = requests[selected-1], requests[selected]
		selected--
		saveRequests()
		printList(v)
		if dv, err := g.View("details"); err == nil {
			printDetails(g, dv)
		}
	}
	return nil
}

func moveRequestDown(g *gocui.Gui, v *gocui.View) error {
	if selected < len(requests)-1 {
		// Swap mit nächstem
		requests[selected], requests[selected+1] = requests[selected+1], requests[selected]
		selected++
		saveRequests()
		printList(v)
		if dv, err := g.View("details"); err == nil {
			printDetails(g, dv)
		}
	}
	return nil
}

func cursorDownDetails(g *gocui.Gui, v *gocui.View) error {
	if inEditPopup {
		return nil
	}
	if detailSelected < 3 {
		detailSelected++
		printDetails(g, v)
	}
	return nil
}

func cursorUpDetails(g *gocui.Gui, v *gocui.View) error {
	if inEditPopup {
		return nil
	}
	if detailSelected > 0 {
		detailSelected--
		printDetails(g, v)
	}
	return nil
}

func openFieldEdit(g *gocui.Gui, v *gocui.View) error {
	if inEditPopup || len(requests) == 0 {
		return nil
	}

	maxX, maxY := g.Size()
	ev, err := g.SetView("fieldEdit", maxX/6, maxY/6, maxX*5/6, maxY*5/6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ev.Title = " Edit Value (Ctrl+S=Save, Esc=Cancel, Ctrl+V=Paste Clipboard) "
		ev.Editable = true
		ev.Wrap = true
		ev.BgColor = gocui.ColorYellow
		ev.FgColor = gocui.ColorBlack
		inEditPopup = true
		g.Cursor = true

		// aktuellen Wert einsetzen
		r := requests[selected]
		var text string
		switch detailSelected {
		case 0:
			text = r.Name
		case 1:
			text = r.Method
		case 2:
			text = r.URL
		case 3:
			text = r.Body
		}
		fmt.Fprint(ev, text)

		// Cursor ans Ende setzen
		ev.SetCursor(len(text), 0)

		// Keybinding für Clipboard einfügen
		g.SetKeybinding("fieldEdit", gocui.KeyCtrlV, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			clip, err := clipboard.ReadAll()
			if err != nil {
				return nil
			}
			v.Clear()
			fmt.Fprint(v, clip)
			return nil
		})

	}

	_, err = g.SetCurrentView("fieldEdit")
	return err
}

func cancelFieldEdit(g *gocui.Gui, v *gocui.View) error {
	log.Println(">>> cancelFieldEdit called")
	g.DeleteView("fieldEdit")
	inEditPopup = false
	g.Cursor = false
	if dv, err := g.View("details"); err == nil {
		g.SetCurrentView("details")
		printDetails(g, dv)
	}
	return nil
}

func saveFieldEdit(g *gocui.Gui, v *gocui.View) error {
	value := strings.TrimSpace(v.Buffer())
	r := &requests[selected]
	switch detailSelected {
	case 0:
		r.Name = value
	case 1:
		r.Method = value
	case 2:
		r.URL = value
	case 3:
		r.Body = value
	}
	saveRequests()
	g.DeleteView("fieldEdit")
	inEditPopup = false
	g.Cursor = false
	if dv, err := g.View("details"); err == nil {
		g.SetCurrentView("details")
		printDetails(g, dv)
	}
	return nil
}

// ---------- List Navigation ----------

func cursorDownList(g *gocui.Gui, v *gocui.View) error {
	if selected < len(requests)-1 {
		selected++
		if dv, err := g.View("details"); err == nil {
			printDetails(g, dv)
		}
	}
	printList(v)
	return nil
}

func cursorUpList(g *gocui.Gui, v *gocui.View) error {
	if selected > 0 {
		selected--
		if dv, err := g.View("details"); err == nil {
			printDetails(g, dv)
		}
	}
	printList(v)
	return nil
}

func openHelp(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetView("help", 10, 5, 70, 20); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v, _ := g.View("help")
		v.BgColor = gocui.ColorYellow
		v.FgColor = gocui.ColorBlack
		v.Title = " Hilfe (Esc zum Schließen) "
		v.Wrap = true

		helpText1 := "  F1          : Hilfe anzeigen"
		helpText2 := "  Arrow Up    : Auswahl nach oben"
		helpText3 := "  Arrow Down  : Auswahl nach unten"
		helpText4 := "  Enter       : Request senden"
		helpText5 := "  Delete      : Request löschen"
		helpText6 := "  PgUp / PgDn : Request verschieben"
		helpText7 := "  e           : Request editieren"
		helpText8 := "  Esc         : Popup schließen / Beenden"

		fmt.Fprint(v, "\n\n")
		fmt.Fprintln(v, helpText1)
		fmt.Fprintln(v, helpText2)
		fmt.Fprintln(v, helpText3)
		fmt.Fprintln(v, helpText4)
		fmt.Fprintln(v, helpText5)
		fmt.Fprintln(v, helpText6)
		fmt.Fprintln(v, helpText7)
		fmt.Fprintln(v, helpText8)
		if _, err := g.SetCurrentView("help"); err != nil {
			return err
		}
	}
	return nil
}

func closeHelp(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("help")
	// Fokus wieder auf Liste setzen
	if _, err := g.View("list"); err == nil {
		g.SetCurrentView("list")
	}
	return nil
}

func sendRequest(g *gocui.Gui, v *gocui.View) error {
	processRequest(g, requests[selected])
	return nil
}

func openResponseView(g *gocui.Gui, content string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("response", 2, 2, maxX-3, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Response (Esc = close) "
		v.Wrap = true
		v.Autoscroll = false // wir scrollen manuell
		v.Editable = false
		v.Clear()
		fmt.Fprint(v, content)

		// --- HIER automatisch nach unten scrollen ---
		lines := strings.Count(content, "\n")
		if lines > maxY { // nur scrollen, wenn mehr Zeilen als Platz
			v.SetOrigin(0, lines-(maxY-5))
		}

		// Keybindings für Scrollen
		g.SetKeybinding("response", gocui.KeyArrowUp, gocui.ModNone, scrollResponseUp)
		g.SetKeybinding("response", gocui.KeyArrowDown, gocui.ModNone, scrollResponseDown)
		g.SetKeybinding("response", gocui.KeyPgup, gocui.ModNone, scrollResponsePgUp)
		g.SetKeybinding("response", gocui.KeyPgdn, gocui.ModNone, scrollResponsePgDn)

		// Schließen mit Esc
		g.SetKeybinding("response", gocui.KeyEsc, gocui.ModNone, closeResponseView)
	}
	_, err := g.SetCurrentView("response")
	return err
}

func scrollResponseUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if oy > 0 {
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

func scrollResponseDown(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	v.SetOrigin(ox, oy+1)
	return nil
}

func scrollResponsePgUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if oy > 10 {
		v.SetOrigin(ox, oy-10)
	} else {
		v.SetOrigin(ox, 0)
	}
	return nil
}

func scrollResponsePgDn(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	v.SetOrigin(ox, oy+10)
	return nil
}

func closeResponseView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("response")
	if _, err := g.View("list"); err == nil {
		g.SetCurrentView("list")
	}
	return nil
}

// ---------- Main ----------

func main() {
	loadRequests()
	if err := run(); err != nil {
		log.Panicln(err)
	}
}

func run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.InputEsc = true // <-- WICHTIG

	g.SetManagerFunc(layout)

	// Keybindings
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, openHelp)

	g.SetKeybinding("help", gocui.KeyEsc, gocui.ModNone, closeHelp)

	g.SetKeybinding("list", gocui.KeyEsc, gocui.ModNone, quit)
	g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, cursorDownList)
	g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, cursorUpList)
	g.SetKeybinding("list", 'e', gocui.ModNone, editRequest)
	g.SetKeybinding("list", gocui.KeyDelete, gocui.ModNone, deleteRequest)
	g.SetKeybinding("list", gocui.KeyPgup, gocui.ModNone, moveRequestUp)
	g.SetKeybinding("list", gocui.KeyPgdn, gocui.ModNone, moveRequestDown)
	g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, sendRequest)

	g.SetKeybinding("details", gocui.KeyArrowDown, gocui.ModNone, cursorDownDetails)
	g.SetKeybinding("details", gocui.KeyArrowUp, gocui.ModNone, cursorUpDetails)
	g.SetKeybinding("details", gocui.KeyEnter, gocui.ModNone, openFieldEdit)
	g.SetKeybinding("details", gocui.KeyEsc, gocui.ModNone, exitEditRequest)

	g.SetKeybinding("fieldEdit", gocui.KeyEsc, gocui.ModNone, cancelFieldEdit)
	g.SetKeybinding("fieldEdit", gocui.KeyCtrlS, gocui.ModNone, saveFieldEdit)

	return g.MainLoop()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func processRequest(g *gocui.Gui, r Request) {
	method := strings.TrimSpace(strings.ToUpper(r.Method))
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH":
		fire_request(g, method, r.URL, r.Body)
	default:
		g.Update(func(g *gocui.Gui) error {
			return openResponseView(g, "UNKNOWN HTTP METHOD")
		})
	}
}

func fire_request(g *gocui.Gui, method string, url string, data string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	showResponse(g, resp)
}

func showResponse(g *gocui.Gui, resp *http.Response) {
	var sb strings.Builder

	if resp.StatusCode == 200 {
		sb.WriteString(fmt.Sprintf("%sResponse status: %d%s\n", green, resp.StatusCode, reset))
	} else {
		sb.WriteString(fmt.Sprintf("%sResponse status: %d%s\n", red, resp.StatusCode, reset))
	}

	sb.WriteString(fmt.Sprintf("%sResponse Headers:%s\n", yellow, reset))
	for key, values := range resp.Header {
		for _, v := range values {
			sb.WriteString(fmt.Sprintf(yellow+"    %s: %s\n", key, v))
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		sb.WriteString(fmt.Sprintf("%s%s%s\n", white, scanner.Text(), reset))
	}
	if err := scanner.Err(); err != nil {
		sb.WriteString(fmt.Sprintf("%sERROR beim Lesen des Bodys: %v%s\n", red, err, reset))
	}

	g.Update(func(g *gocui.Gui) error {
		return openResponseView(g, sb.String())
	})
}
