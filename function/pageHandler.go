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

func Scattegories(w http.ResponseWriter, r *http.Request, letter string) {
	template, err := template.ParseFiles("./pages/scattegories.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, letter)
}

func ScattegoriesVerification(w http.ResponseWriter, r *http.Request, data Question) {
	template, err := template.ParseFiles("./pages/sctVerification.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

func Waiting(w http.ResponseWriter, r *http.Request) {
	userid := GetCoockie(w, r, "userId")
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, userid)
}

func DeafTest(w http.ResponseWriter, r *http.Request, currentMusic Music) {
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, currentMusic)
}
