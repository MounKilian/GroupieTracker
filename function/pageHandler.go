package groupieTracker

import (
	"database/sql"
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
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if userid == 1 {
		values := [2]int{53, 1}
		addPlayer(db, values)
		template, err := template.ParseFiles("./pages/waitingInvit.html")
		if err != nil {
			log.Fatal(err)
		}
		template.Execute(w, userid)
	} else {
		value := [4]string{"3", "6", "TestGame", "3"}
		createNewRoom(db, value)
		template, err := template.ParseFiles("./pages/waiting.html")
		if err != nil {
			log.Fatal(err)
		}
		template.Execute(w, userid)
	}
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/verification", http.StatusFound)
}
