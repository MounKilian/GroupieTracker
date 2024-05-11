package groupieTracker

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
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
	code := GetCoockieCode(w, r, "code")
	var newInfo Info
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if code == "" {
		user := GetUserById(db, userid)
		IDroom, err := GetRoomByName(db, responseJoin)
		if IDroom == 0 {
			http.Redirect(w, r, "/room", http.StatusFound)
		} else {
			if err != nil {
				log.Fatal(err)
			}
			values := [2]int{IDroom, userid}
			addPlayer(db, values)
			users, err := getUsersInRoom(db, responseJoin)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range users {
				newInfo.Pseudo = append(newInfo.Pseudo, i)
			}
			SetCookieCode(w, user, responseJoin)
			newInfo = Info{responseJoin, newInfo.Pseudo}
			room.broadcastMessage("newUser")
		}
	} else {
		users, err := getUsersInRoom(db, code)
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range users {
			newInfo.Pseudo = append(newInfo.Pseudo, i)
		}
		newInfo = Info{code, newInfo.Pseudo}
	}
	template, err := template.ParseFiles("./pages/waitingInvit.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, newInfo)
}

func Waiting(w http.ResponseWriter, r *http.Request, room *Room) {
	userid := GetCoockie(w, r, "userId")
	code := GetCoockieCode(w, r, "code")
	db, err := sql.Open("sqlite3", "BDD.db")
	var newInfo Info
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if code == "" {
		infos.Pseudo = []string{}
		newInfo.Pseudo = []string{}
		user := GetUserById(db, userid)
		pseudo := user.pseudo
		access := RandomString()
		SetCookieCode(w, user, access)
		var value [4]string
		if game == "scattegories" {
			value = [4]string{strconv.Itoa(userid), "6", access, "3"}
		} else if game == "blindTest" {
			value = [4]string{strconv.Itoa(userid), "6", access, "1"}
		} else {
			value = [4]string{strconv.Itoa(userid), "6", access, "2"}
		}
		createNewRoom(db, value)
		infos.Pseudo = append(infos.Pseudo, pseudo)
		infos = Info{access, infos.Pseudo}
		newInfo.Pseudo = append(infos.Pseudo, pseudo)
		newInfo = Info{access, infos.Pseudo}
	} else {
		newInfo.Pseudo = []string{}
		users, err := getUsersInRoom(db, code)
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range users {
			newInfo.Pseudo = append(newInfo.Pseudo, i)
		}
		newInfo = Info{code, newInfo.Pseudo}
		log.Println(infos.Pseudo)
	}
	template, err := template.ParseFiles("./pages/waiting.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, newInfo)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/verification", http.StatusFound)
}

func Win(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userId := GetCoockie(w, r, "userId")
	currentRoom := GetCurrentRoomUser(db, userId)
	playersScore := GetUsersScoreInRoom(db, currentRoom)

	template, err := template.ParseFiles("./pages/win.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, playersScore)
}
