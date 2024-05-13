package groupieTracker

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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

type Note struct {
	ID int

	Artist struct {
		True  int
		False int
		Same  int
	}

	Album struct {
		True  int
		False int
		Same  int
	}

	Groupe struct {
		True  int
		False int
		Same  int
	}

	Instrument struct {
		True  int
		False int
		Same  int
	}

	Featuring struct {
		True  int
		False int
		Same  int
	}
}

var Reponse []Note

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
			currentId, err := strconv.Atoi(ExtractSuffix(key))
			if err != nil {
				log.Fatal(err)
			}
			var found bool
			for index, note := range Reponse {
				id := note.ID
				if currentId == id {
					modifyCategoryField(ExtractSuffix(key), "True", 1, &Reponse[index])
					found = true
				}
			}
			if found != true {
				newNote := Note{ID: currentId}
				Reponse = append(Reponse, newNote)
				modifyCategoryField(ExtractSuffix(key), "True", 1, &Reponse[len(Reponse)-1])
			}
			break
		case "same":
			currentId, err := strconv.Atoi(ExtractSuffix(key))
			if err != nil {
				log.Fatal(err)
			}
			var found bool
			for index, note := range Reponse {
				id := note.ID
				if currentId == id {
					modifyCategoryField(ExtractSuffix(key), "Same", 1, &Reponse[index])
					found = true
				}
			}
			if found != true {
				newNote := Note{ID: currentId}
				Reponse = append(Reponse, newNote)
				modifyCategoryField(ExtractSuffix(key), "Same", 1, &Reponse[len(Reponse)-1])
			}
			break
		case "false":
			currentId, err := strconv.Atoi(ExtractSuffix(key))
			if err != nil {
				log.Fatal(err)
			}
			var found bool
			for index, note := range Reponse {
				id := note.ID
				if currentId == id {
					modifyCategoryField(ExtractSuffix(key), "False", 1, &Reponse[index])
					found = true
				}
			}
			if found != true {
				newNote := Note{ID: currentId}
				Reponse = append(Reponse, newNote)
				modifyCategoryField(ExtractSuffix(key), "False", 1, &Reponse[len(Reponse)-1])
			}
			break
		}
	}

	http.Redirect(w, r, "/win", http.StatusFound)
}

func modifyCategoryField(categoryName string, fieldName string, value int, note *Note) {
	categoryFields := strings.Split(categoryName, ".")
	currentField := reflect.ValueOf(note).Elem()

	for i := 0; i < len(categoryFields)-1; i++ {
		field := currentField.FieldByName(categoryFields[i])
		if !field.IsValid() {
			fmt.Printf("Field %s not found in Note struct\n", categoryFields[i])
			return
		}
		currentField = field.Elem()
	}

	field := currentField.FieldByName(categoryFields[len(categoryFields)-1])
	if !field.IsValid() {
		fmt.Printf("Field %s not found in Note struct\n", categoryFields[len(categoryFields)-1])
		return
	}

	fieldType := field.Type()
	valueType := reflect.TypeOf(value)

	if fieldType != valueType {
		fmt.Printf("Type mismatch between category field %s and value %d\n", categoryName, value)
		return
	}

	// Get the nested field by name
	nestedField := field.FieldByName(fieldName)
	if !nestedField.IsValid() {
		fmt.Printf("Field %s not found in category field %s\n", fieldName, categoryFields[len(categoryFields)-1])
		return
	}

	currentValue := nestedField.Int()
	newValue := currentValue + int64(value)
	nestedField.SetInt(newValue)
}
