package groupieTracker

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// main pages
func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Waiting(w http.ResponseWriter, r *http.Request) {
	userid := GetCoockie(w, r, "userId")
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, userid)
}

// Scattegories pages
func Scattegories(w http.ResponseWriter, r *http.Request, letter string) {
	template, err := template.ParseFiles("./pages/scategories/scattegories.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, letter)
}

func ScattegoriesVerification(w http.ResponseWriter, r *http.Request, data Question) {
	template, err := template.ParseFiles("./pages/scategories/verification.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

// Deaftest pages
func DeafTest(w http.ResponseWriter, r *http.Request, Deaftest Deaftest) {
	template, err := template.ParseFiles("./pages/deaftest.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, Deaftest.currentMusic.lyrics)
}

func Win(w http.ResponseWriter, r *http.Request, currentMusic Music) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userId := GetCoockie(w, r, "userId")
	currentRoom := GetCurrentRoomUser(db, userId)
	score := GetPlayerScore(db, currentRoom, userId)

	template, err := template.ParseFiles("./pages/win.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, score)
}
