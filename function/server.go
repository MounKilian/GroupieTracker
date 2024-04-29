package groupieTracker

import (
	"fmt"
	"net/http"
)

var questions = []Question{}

func Server() {
	var currentUser User
	var letter string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r)
	})
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		currentUser = Formulaire(w, r)
		fmt.Println(currentUser)
	})
	http.HandleFunc("/scattegories", func(w http.ResponseWriter, r *http.Request) {
		letter = selectRandomLetter()
		Scattegories(w, r, letter)
	})
	http.HandleFunc("/scattegoriesChecker", func(w http.ResponseWriter, r *http.Request) {
		response := ScattegoriesForm(w, r)
		questions = append(questions, response)
		fmt.Println(questions)
	})
	http.HandleFunc("/verification", func(w http.ResponseWriter, r *http.Request) {
		ScattegoriesVerification(w, r, questions[0])
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}
