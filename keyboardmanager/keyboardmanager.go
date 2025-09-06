package keyboardmanager

import (
	"os"
)

// Key repräsentiert eine Taste
type Key string

const (
	KeyUnknown    Key = "Unknown"
	KeyEnter      Key = "Enter"
	KeyBackspace  Key = "Backspace"
	KeyDelete     Key = "Delete"
	KeyArrowUp    Key = "ArrowUp"
	KeyArrowDown  Key = "ArrowDown"
	KeyArrowLeft  Key = "ArrowLeft"
	KeyArrowRight Key = "ArrowRight"
	KeyEscape     Key = "Escape"
	KeyF1         Key = "F1"
	KeyF2         Key = "F2"
	KeyF3         Key = "F3"
	KeyF4         Key = "F4"
	KeyF5         Key = "F5"
	KeyF6         Key = "F6"
	KeyF7         Key = "F7"
	KeyF8         Key = "F8"
	KeyF9         Key = "F9"
	KeyF10        Key = "F10"
	KeyF11        Key = "F11"
	KeyF12        Key = "F12"
)

// ReadKey liest eine Taste vom Terminal und gibt den Key zurück
func ReadKey() Key {
	buf := make([]byte, 5)
	n, _ := os.Stdin.Read(buf)

	if n == 0 {
		return KeyUnknown
	}

	// Normale Zeichen
	if n == 1 {
		switch buf[0] {
		case 13, 10:
			return KeyEnter
		case 27:
			return KeyEscape // ← ESC direkt erkannt
		case 127:
			return KeyBackspace
		default:
			return Key(buf[0:1])
		}
	}

	// Escape-Sequenzen
	if buf[0] == 27 {
		// Pfeiltasten / Delete ESC [
		if n >= 3 && buf[1] == 91 {
			switch buf[2] {
			case 65:
				return KeyArrowUp
			case 66:
				return KeyArrowDown
			case 67:
				return KeyArrowRight
			case 68:
				return KeyArrowLeft
			case 51: // Delete
				if n >= 4 && buf[3] == 126 {
					return KeyDelete
				}
			}
		}
		// F1–F4 ESC O
		if n >= 3 && buf[1] == 79 {
			switch buf[2] {
			case 80:
				return KeyF1
			case 81:
				return KeyF2
			case 82:
				return KeyF3
			case 83:
				return KeyF4
			}
		}
		// F5–F12 ESC [ ..
		if n >= 5 && buf[1] == 91 {
			switch {
			case buf[2] == 49 && buf[3] == 53 && buf[4] == 126:
				return KeyF5
			case buf[2] == 49 && buf[3] == 55 && buf[4] == 126:
				return KeyF6
			case buf[2] == 49 && buf[3] == 56 && buf[4] == 126:
				return KeyF7
			case buf[2] == 49 && buf[3] == 57 && buf[4] == 126:
				return KeyF8
			case buf[2] == 50 && buf[3] == 48 && buf[4] == 126:
				return KeyF9
			case buf[2] == 50 && buf[3] == 49 && buf[4] == 126:
				return KeyF10
			case buf[2] == 50 && buf[3] == 51 && buf[4] == 126:
				return KeyF11
			case buf[2] == 50 && buf[3] == 52 && buf[4] == 126:
				return KeyF12
			}
		}
	}

	return KeyUnknown
}
