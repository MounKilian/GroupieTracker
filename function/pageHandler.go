package groupieTracker

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func RoomStart(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/room.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func WaitingInvit(w http.ResponseWriter, r *http.Request) {
	responseJoin := r.FormValue("join")
	userid := GetCoockie(w, r, "userId")
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	IDroom, err := GetRoomByName(db, responseJoin)
	if err != nil {
		log.Fatal(err)
	}
	values := [2]int{IDroom, userid}
	addPlayer(db, values)
	template, err := template.ParseFiles("./pages/waitingInvit.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, responseJoin)
}

func Waiting(w http.ResponseWriter, r *http.Request) {
	userid := GetCoockie(w, r, "userId")
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	access := RandomString()
	value := [4]string{strconv.Itoa(userid), "6", access, "3"}
	createNewRoom(db, value)
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, access)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/verification", http.StatusFound)
}
