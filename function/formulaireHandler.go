package groupieTracker

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Question struct {
	Id                int
	Username          string
	Artist            string
	Album             string
	GroupeDeMusic     string
	InstrumentDeMusic string
	Featuring         string
}

func ScattegoriesForm(w http.ResponseWriter, r *http.Request) Question {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	userId := GetCoockie(w, r, "userId")

	var response Question
	response.Id = userId
	response.Username = GetUserById(db, userId).pseudo
	response.Artist = r.FormValue("artist")
	response.Album = r.FormValue("album")
	response.GroupeDeMusic = r.FormValue("groupe-de-music")
	response.InstrumentDeMusic = r.FormValue("instrument-de-music")
	response.Featuring = r.FormValue("featuring")

	http.Redirect(w, r, "/verification", http.StatusFound)
	return response
}

func Formulaire(w http.ResponseWriter, r *http.Request) {
	DeleteCookies(w, r)
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var user User
	for key := range r.Form {
		if key == "pseudo-sign" {
			pseudo := r.FormValue("pseudo-sign")
			email := r.FormValue("email-sign")
			password := r.FormValue("password-sign")
			verifyPassword := r.FormValue("verify-password-sign")
			passwordEncrypt := Encrypt(password)
			valueCreate := [3]string{pseudo, email, passwordEncrypt}

			if verifyPassword != password && VerifyPassword(password) {
				log.Println("ERROR : Incorrect password")
			} else {
				createUser(db, valueCreate)
				valueConnect := [2]string{pseudo, passwordEncrypt}
				user = connectUser(db, valueConnect)
				if user.pseudo == "" {
					log.Println("ERROR : Wrong connection information")
				} else {
					SetCookie(w, user)
					http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				}
			}
		} else if key == "connect-log" {
			password := r.FormValue("password-log")
			passwordEncrypt := Encrypt(password)
			connect := r.FormValue("connect-log")
			valueConnect := [2]string{connect, passwordEncrypt}

			user = connectUser(db, valueConnect)
			if user.pseudo == "" {
				log.Println("ERROR : Wrong connection information")
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			} else {
				SetCookie(w, user)
				http.Redirect(w, r, "/landingPage", http.StatusFound)
			}
		}
	}
}

func DeafForm(w http.ResponseWriter, r *http.Request, Deaftest *Deaftest) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	musicName := r.FormValue("music-name")
	if MatchTitle(Deaftest.currentMusic.name, musicName) {
		userId := GetCoockie(w, r, "userId")
		currentRoom := GetCurrentRoomUser(db, userId)
		UpdatePlayerScore(db, currentRoom, userId, 10)
	}
	if Deaftest.currentPlay == Deaftest.nbSong {
		http.Redirect(w, r, "/win", http.StatusFound)
	}

	http.Redirect(w, r, "/deaftestround", http.StatusFound)
}

func BlindForm(w http.ResponseWriter, r *http.Request, Blindtest *Blindtest) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	musicName := r.FormValue("music-name")
	if MatchTitle(Blindtest.currentBtest.name, musicName) {
		userId := GetCoockie(w, r, "userId")
		currentRoom := GetCurrentRoomUser(db, userId)
		UpdatePlayerScore(db, currentRoom, userId, 10)
	}
	if Blindtest.currentPlay == Blindtest.nbSong {
		http.Redirect(w, r, "/win", http.StatusFound)
	}
	http.Redirect(w, r, "/blindtestround", http.StatusFound)
}

func WaitingForm(w http.ResponseWriter, r *http.Request) (string, int) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	var playlistID string
	gender := r.FormValue("gender")
	nbSong, err := strconv.Atoi(r.FormValue("nb-song"))
	if err != nil {
		log.Fatal(err)
	}

	switch gender {
	case "rock":
		playlistID = "2xF91FPK8mitD7ND6iAt5j"
		break
	case "pop":
		playlistID = "4SSiAXhcLdrGSCGpL1B8wG"
		break
	case "us-rap":
		playlistID = "5VvIIqAZBwjQ7hTqZbcjKr"
		break
	case "normal":
		playlistID = "5tYg6pvAiwa3taoNAG3HzC"
		break
	default:
		playlistID = "5tYg6pvAiwa3taoNAG3HzC"
		break
	}

	http.Redirect(w, r, "/deaftest", http.StatusFound)
	return playlistID, nbSong
}

func WaitingFormBT(w http.ResponseWriter, r *http.Request) (string, int) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	var playlistID string
	gender := r.FormValue("gender")
	nbSong, err := strconv.Atoi(r.FormValue("nb-song"))
	if err != nil {
		log.Fatal(err)
	}

	switch gender {
	case "rock":
		playlistID = "2xF91FPK8mitD7ND6iAt5j"
		break
	case "pop":
		playlistID = "4SSiAXhcLdrGSCGpL1B8wG"
		break
	case "us-rap":
		playlistID = "5VvIIqAZBwjQ7hTqZbcjKr"
		break
	case "normal":
		playlistID = "5tYg6pvAiwa3taoNAG3HzC"
		break
	default:
		playlistID = "5tYg6pvAiwa3taoNAG3HzC"
		break
	}

	http.Redirect(w, r, "/blindtest", http.StatusFound)
	return playlistID, nbSong
}

func ScattegoriesVerificationChecker(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	for key := range r.Form {
		switch r.FormValue(key) {
		case "true":
			userId, err := strconv.Atoi(ExtractSuffix(key))
			if err != nil {
				log.Fatal(err)
			}
			currentRoom := GetCurrentRoomUser(db, userId)
			UpdatePlayerScore(db, currentRoom, userId, 3)
			break
		case "same":
			userId, err := strconv.Atoi(ExtractSuffix(key))
			if err != nil {
				log.Fatal(err)
			}
			currentRoom := GetCurrentRoomUser(db, userId)
			UpdatePlayerScore(db, currentRoom, userId, 1)
			break
		}
	}

	http.Redirect(w, r, "/win", http.StatusFound)
}
