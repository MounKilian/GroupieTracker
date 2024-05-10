package groupieTracker

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id       int
	pseudo   string
	email    string
	password string
}

func selectFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

// Find occurence for a specifique user
func checkUser(db *sql.DB, value [3]string) int {
	var nbAccount int

	query := "SELECT COUNT(*) FROM USER WHERE pseudo = ? OR email = ?"
	err := db.QueryRow(query, value[0], value[1]).Scan(&nbAccount)
	if err != nil {
		log.Fatal(err)
	}
	return nbAccount
}

// create a new user in database
func createUser(db *sql.DB, value [3]string) {
	if checkUser(db, value) == 0 {
		insertQuery := "INSERT INTO USER (id, pseudo, email, password) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, nil, value[0], value[1], value[2])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Return user information as User
func connectUser(db *sql.DB, value [2]string) User {
	var u User
	db.QueryRow("SELECT * FROM `USER` WHERE pseudo = ? OR email = ? AND password = ?", value[0], value[0], value[1]).Scan(&u.id, &u.pseudo, &u.email, &u.password)
	return u
}

// Print specifique rows
func displayUserRows(rows *sql.Rows) {
	for rows.Next() {
		var u User
		err := rows.Scan(&u.id, &u.pseudo, &u.email, &u.password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u)
	}
}

// Reset User table
func resetUserTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS `USER`")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE USER (id INTEGER PRIMARY KEY, pseudo TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

// Create a new room and add creator to the new room
func createNewRoom(db *sql.DB, value [4]string) {
	//Convert to int
	nbMaxPlayer, _ := strconv.Atoi(value[1])
	created_by, _ := strconv.Atoi(value[0])
	id_game, _ := strconv.Atoi(value[3])
	if nbMaxPlayer > 6 {
		log.Println("to much player")
	} else {
		insertQuery := "INSERT INTO ROOMS (id, created_by, max_player, name, id_game) VALUES (?, ?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, nil, created_by, nbMaxPlayer, value[2], id_game)
		if err != nil {
			log.Fatal(err)
		}

		currentUserValue := [2]int{getRoomCreator(db, created_by), created_by}
		AddPlayer(db, currentUserValue)
	}
}

// add player in a specifique room
func AddPlayer(db *sql.DB, value [2]int) {
	if checkNbPlayer(db, value[0]) > getMaxPlayer(db, value[0]) {
		log.Println("Party already full")
	} else {
		insertQuery := "INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)"
		_, err := db.Exec(insertQuery, value[0], value[1], 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Increments player's score
func UpdatePlayerScore(db *sql.DB, id_room int, id_user int, score int) {
	insertQuery := "UPDATE ROOM_USERS SET score = score + ? WHERE id_room = ? AND id_user = ?"
	_, err := db.Exec(insertQuery, score, id_room, id_user)
	if err != nil {
		log.Fatal(err)
	}
}

// return score for a specifique room user
func GetPlayerScore(db *sql.DB, id_room int, id_user int) int {
	var score int

	query := "SELECT score FROM ROOM_USERS WHERE id_user = ? and id_room = ?"
	err := db.QueryRow(query, id_user, id_room).Scan(&score)
	if err != nil {
		log.Fatal(err)
	}
	return score
}

// Check if room is not full
func checkNbPlayer(db *sql.DB, id_room int) int {
	var nbPLayer int
	query := "SELECT COUNT(*) FROM ROOM_USERS WHERE id_room = ?"
	err := db.QueryRow(query, id_room).Scan(&nbPLayer)
	if err != nil {
		log.Fatal(err)
	}
	return nbPLayer
}

// return max player of room
func getMaxPlayer(db *sql.DB, id_room int) int {
	var nbMaxPlayer int

	query := "SELECT `max_player` FROM ROOMS WHERE id = ?"
	err := db.QueryRow(query, id_room).Scan(&nbMaxPlayer)
	if err != nil {
		log.Fatal(err)
	}
	return nbMaxPlayer
}

// return the creator of a specifique room
func getRoomCreator(db *sql.DB, id_user int) int {
	var idRoom int

	query := "SELECT MAX(id) FROM ROOMS WHERE created_by = ?"
	err := db.QueryRow(query, id_user).Scan(&idRoom)
	if err != nil {
		log.Fatal(err)
	}
	return idRoom
}

func GetCrteatedPlayer(db *sql.DB, id_room int) int {
	var idCreatedPlayer int

	query := "SELECT created_by FROM ROOMS WHERE id = ?"
	err := db.QueryRow(query, id_room).Scan(&idCreatedPlayer)
	if err != nil {
		log.Fatal(err)
	}
	return idCreatedPlayer
}

// return user information as User for a specifique id
func GetUserById(db *sql.DB, id int) User {
	var u User
	db.QueryRow("SELECT * FROM `USER` WHERE id = ?", id).Scan(&u.id, &u.pseudo, &u.email, &u.password)
	return u
}

// get last room (current room) for a specifique user
func GetCurrentRoomUser(db *sql.DB, idUser int) int {
	var idRoom int

	query := "SELECT id_room FROM ROOM_USERS WHERE id_user = ? ORDER BY id_room DESC LIMIT 1"
	err := db.QueryRow(query, idUser).Scan(&idRoom)
	if err != nil {
		log.Fatal(err)
	}
	return idRoom
}

func GetRoomByName(db *sql.DB, roomName string) (int, error) {
	var idRoom int

	query := "SELECT id FROM ROOMS WHERE name =?"
	err := db.QueryRow(query, roomName).Scan(&idRoom)
	if err != nil {
		return 0, err
	}
	return idRoom, nil
}

func getUsersInRoom(db *sql.DB, roomName string) ([]string, error) {
	var users []string
	query := "SELECT u.pseudo FROM ROOM_USERS ru JOIN USER u ON ru.id_user = u.id WHERE ru.id_room = (SELECT id FROM ROOMS WHERE name = $1)"
	rows, err := db.Query(query, roomName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pseudo string
		err := rows.Scan(&pseudo)
		if err != nil {
			return nil, err
		}
		users = append(users, pseudo)
	}
	return users, nil
}
