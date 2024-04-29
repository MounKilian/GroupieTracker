package groupieTracker

import (
	"database/sql"
	"log"
	"net/http"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

type Question struct {
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
	response.Username = GetUserById(db, userId).pseudo
	response.Artist = r.FormValue("artist")
	response.Album = r.FormValue("album")
	response.GroupeDeMusic = r.FormValue("groupe-de-music")
	response.InstrumentDeMusic = r.FormValue("instrument-de-music")
	response.Featuring = r.FormValue("featuring")

	http.Redirect(w, r, "/verification", http.StatusFound)
	return response
}

func Formulaire(w http.ResponseWriter, r *http.Request) User {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var user User
	for key, _ := range r.Form {
		if key == "pseudo-sign" {
			pseudo := r.FormValue("pseudo-sign")
			email := r.FormValue("email-sign")
			password := r.FormValue("password-sign")
			verifyPassword := r.FormValue("verify-password-sign")
			passwordEncrypt := Encrypt(password)
			valueCreate := [3]string{pseudo, email, passwordEncrypt}

			if verifyPassword != password && VerifyPassword(password) {
				log.Println("ERROR : Incorrect password")
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			} else {
				createUser(db, valueCreate)
				valueConnect := [2]string{pseudo, passwordEncrypt}
				user = connectUser(db, valueConnect)
				if user.pseudo == "" {
					log.Println("ERROR : Wrong connection information")
				} else {
					http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
					return user
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
			} else {
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				return user
			}
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	return user
}

func VerifyPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}
