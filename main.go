// static-files.go
package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templatesPath = "templates/*.html"
var tpl *template.Template

func main() {
	var err error
	tpl, err = template.ParseGlob(templatesPath)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/about", AboutHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "about", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "home", nil)
	if err != nil {
		log.Fatal(err)
	}

}
