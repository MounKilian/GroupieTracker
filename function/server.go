package groupieTracker

import (
	"fmt"
	"net/http"
)

func Server() {
	var currentUser User

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		errorMessage := Error{}
		Home(w, r, errorMessage)
	})
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		currentUser = Formulaire(w, r)
		fmt.Println(currentUser)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}
