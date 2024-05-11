package groupieTracker

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

// Blindtest pages
func BlindTest(w http.ResponseWriter, r *http.Request, Blindtest Blindtest) {
	template, err := template.ParseFiles("./pages/blindtest.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, Blindtest.currentBtest.PreviewURL)
}
