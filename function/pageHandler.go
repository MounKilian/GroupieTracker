package groupieTracker

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Info struct {
	Code   string
	Pseudo []string
}

var game string
var infos Info
var refresh = true

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

func LandingPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/landingPage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func ScattegoriesVerification(w http.ResponseWriter, r *http.Request, data Question) {
	template, err := template.ParseFiles("./pages/sctVerification.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

func RoomStart(w http.ResponseWriter, r *http.Request) {
	game = r.FormValue("game-value")
	log.Println(game)
	template, err := template.ParseFiles("./pages/room.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func WaitingInvit(w http.ResponseWriter, r *http.Request, room *Room) {
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
	user := GetUserById(db, userid)
	pseudo := user.pseudo
	for _, i := range infos.Pseudo {
		if i == pseudo {
			refresh = false
		}
	}
	if refresh {
		infos.Pseudo = append(infos.Pseudo, pseudo)
		infos = Info{responseJoin, infos.Pseudo}
		room.broadcastMessage("newUser")
	}
	refresh = true
	template, err := template.ParseFiles("./pages/waitingInvit.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, infos)
}

func Waiting(w http.ResponseWriter, r *http.Request, room *Room) {
	userid := GetCoockie(w, r, "userId")
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	access := RandomString()
	url.QueryEscape(access)
	var value [4]string
	if game == "scattegories" {
		value = [4]string{strconv.Itoa(userid), "6", access, "3"}
	} else if game == "blindTest" {
		value = [4]string{strconv.Itoa(userid), "6", access, "1"}
	} else {
		value = [4]string{strconv.Itoa(userid), "6", access, "2"}
	}
	createNewRoom(db, value)
	user := GetUserById(db, userid)
	pseudo := user.pseudo
	for _, i := range infos.Pseudo {
		if i == pseudo {
			refresh = false
		}
	}
	if refresh {
		infos.Pseudo = append(infos.Pseudo, pseudo)
		infos = Info{access, infos.Pseudo}
		room.broadcastMessage("newUser")
	}
	refresh = true
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, infos)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/verification", http.StatusFound)
}
