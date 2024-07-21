package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type AsciiArtData struct {
	Text     string
	AsciiArt string
	Banner   string
}

// AsciiArtHandler handles POST requests, process input text,
// generate ASCII art and render an HTML template with result.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Get input text and banner from request.
		text := r.FormValue("text")
		banner := r.FormValue("banner")

		// Validate input
		if text == "" || banner == "" {
			http.Error(w, "Error 400: Bad request", http.StatusBadRequest)
			return
		}

		// Read the banner file and generate ASCII art
		lines, err := readBanner(banner)
		if err != nil {
			http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
			return
		}

		// Split the input text into multiple lines
		textLines := strings.Split(text, "\r\n")

		// Process the input text and generate ASII art.
		var asciiArtBuffer bytes.Buffer
		for _, words := range textLines {
			for i := 0; i < 8; i++ {
				for _, char := range words {
					if !(char >= 32 && char <= 126) {
						http.Error(w, "Error 400: Bad request", http.StatusBadRequest)
						return
					}
					asciiArtBuffer.WriteString(lines[int(char-' ')*9+1+i] + " ")
				}
				asciiArtBuffer.WriteString("\n")
			}
			asciiArtBuffer.WriteString("\n")
		}

		asciiArt := asciiArtBuffer.String()

		// Render the ASCII art template
		data := AsciiArtData{Text: text, AsciiArt: asciiArt, Banner: banner}
		tmpl, err := template.ParseFiles("../templates/asciiart.html")
		if err != nil {
			http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Error 405: Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ReadBanner reads banner file content and returns splited content as a slice of strings.
func readBanner(banner string) ([]string, error) {
	path := "../banners/" + banner
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Check file's integrity.
	hashFile := checkSum(data)
	hashShadow := "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73"
	hashStandard := "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"
	hashThinkertoy := "092d0cde973bfbb02522f18e00e8612e269f53bac358bb06f060a44abd0dbc52"

	if hashFile != hashShadow && hashFile != hashStandard && hashFile != hashThinkertoy {
		return nil, fmt.Errorf("file corrupted")
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

// CheckSum calculates the SHA-256 checksum of a given byte slice and returns it as a hexadecimal string.
func checkSum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
