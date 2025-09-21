# HTTP Request Manager (TUI)

Ein kleines **Terminal-UI**-Tool (basierend auf [gocui](https://github.com/jroimartin/gocui)),  
um vordefinierte HTTP-Requests zu verwalten und direkt auszuführen.

---

## ✨ Features

- 📂 Requests werden in einer JSON-Datei gespeichert (`requests.json`)
- 📝 CRUD-Operationen auf Requests:
  - Hinzufügen, Bearbeiten, Löschen, Verschieben
- 📡 HTTP-Methoden unterstützt: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`
- 📜 Response wird in einer **scrollbaren Ansicht** angezeigt
- 🎨 Farbiges TUI mit Navigation per Tastatur

---

## 🖥️ Screenshots

*(Platzhalter, hier kannst du später Screenshots einfügen)*

---

## ⌨️ Tastenkürzel

**Allgemein**
- `F1` – Hilfe anzeigen
- `Esc` – Popup schließen / Programm beenden
- `Ctrl+C` – Programm beenden

**Liste**
- `↑ / ↓` – Auswahl bewegen
- `Enter` – Request senden
- `Delete` – Request löschen
- `PgUp / PgDn` – Request verschieben
- `e` – Request bearbeiten

**Details**
- `↑ / ↓` – Feld auswählen
- `Enter` – Feld editieren
- `Esc` – zurück zur Liste

**Response-View**
- `↑ / ↓` – scrollen
- `PgUp / PgDn` – schneller scrollen
- `Esc` – zurück zum Menü

---

## 🚀 Installation & Start

```bash
# Repository klonen
git clone https://github.com/<dein-user>/<repo-name>.git
cd <repo-name>

# Build
go build -o http-tui

# Start
./http-tui
