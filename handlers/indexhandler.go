package handlers

import (
	"html/template"
	"net/http"
	// "ascii/han"
)

// IndexHandler renders an HTML template for the root path.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Handle non-root paths
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		ErrorHandler(w, "page not found", http.StatusNotFound)
		return
	}

	// Handle only Get method for the root path
	if r.Method != http.MethodGet {
		// http.Error(w, "Error 405: Method not allowed", http.StatusMethodNotAllowed)
		ErrorHandler(w, "Error 405: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render the template for the root path
	tmpl, err1 := (template.ParseFiles("../templates/home.html"))
	if err1 != nil {
		// http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
		ErrorHandler(w, "Error 500: Internal server error", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, nil)
	if err != nil {
		// http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
		ErrorHandler(w, "Error 500: Internal server error", http.StatusInternalServerError)
	}
}
