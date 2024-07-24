package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, str string, code int) {
	varerr, err := template.ParseFiles("../templates/error.html")
	if err != nil {
		http.Error(w, "Error 500: Internal server error", http.StatusInternalServerError)
		log.Fatal("Error", err)
		return
	}
	Person := struct {
		Code int
		Str  string
	}{
		Code: code,
		Str:  str,
	}
	w.WriteHeader(code)
	err = varerr.Execute(w, Person)
	if err != nil {
		log.Fatal("Error", err)
	}
}
