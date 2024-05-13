package groupieTracker

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var questions = []Question{}
var questionsMap map[string][]Question

var state = false
var players = 0
var playersBT = 0
var playersSC = 0
var last = 0
var lastBT = 0

type Blindtest struct {
	gender       string
	nbSong       int
	currentBtest Btest
	currentPlay  int
	nbPlayer     int
	finish       bool
	Time         string
}

type Deaftest struct {
	gender       string
	nbSong       int
	currentMusic Music
	currentPlay  int
	nbPlayer     int
	finish       bool
	Time         string
}

type Room struct {
	id             string
	clients        map[*websocket.Conn]bool
	register       chan *websocket.Conn
	unregister     chan *websocket.Conn
	mu             sync.RWMutex
	Letter         string
	game           string
	Time           string
	nbrsMaxPlayers int
	nbPlayer       int
	finish         bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewRoom() *Room {
	return &Room{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		finish:     false,
	}
}

func NewDeaf() *Deaftest {
	return &Deaftest{
		gender:      "",
		nbSong:      0,
		currentPlay: 0,
		nbPlayer:    0,
		finish:      false,
	}
}

func NewBlind() *Blindtest {
	return &Blindtest{
		gender:      "",
		nbSong:      0,
		currentPlay: 0,
		finish:      false,
	}
}

func (room *Room) Start() {
	for {
		select {
		case conn := <-room.register:
			room.clients[conn] = true
			log.Println("Client connected")
		case conn := <-room.unregister:
			delete(room.clients, conn)
			log.Println("Client disconnected")
		}
	}
}

func (room *Room) broadcastMessage(message string) {
	room.mu.Lock()
	for conn := range room.clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(room.clients, conn)
		}
	}
	defer room.mu.Unlock()
}

func Server() {
	room := NewRoom()
	go room.Start()
	Deaftest := NewDeaf()
	Blindtest := NewBlind()
	questionsMap = make(map[string][]Question)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, "")
	})
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		Formulaire(w, r)
	})
	http.HandleFunc("/landingPage", func(w http.ResponseWriter, r *http.Request) {
		LandingPage(w, r)
	})
	http.HandleFunc("/scattegories", func(w http.ResponseWriter, r *http.Request) {
		Scattegories(w, r, room)
	})
	http.HandleFunc("/scattegoriesChecker", func(w http.ResponseWriter, r *http.Request) {
		room.broadcastMessage("refresh")
		buttonValue := r.FormValue("button-value")
		code := GetCoockieCode(w, r, "code")
		userid := GetCoockie(w, r, "userId")
		questions := questionsMap[code]
		if strconv.Itoa(userid) == buttonValue {
			room.broadcastMessage("end_" + code + strconv.Itoa(userid))
			response := ScattegoriesForm(w, r)
			questions = append(questions, response)
			questionsMap[code] = questions
		} else {
			response := ScattegoriesForm(w, r)
			questions = append(questions, response)
			questionsMap[code] = questions
		}
	})
	http.HandleFunc("/verification", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		questions := questionsMap[code]
		ScattegoriesVerification(w, r, questions)
	})
	http.HandleFunc("/verificationChecker", func(w http.ResponseWriter, r *http.Request) {
		ScattegoriesVerificationChecker(w, r)
	})
	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		nbPlayers := r.FormValue("nb-player")
		log.Println(nbPlayers)
		room.nbrsMaxPlayers, _ = strconv.Atoi(nbPlayers)
		Waiting(w, r, room)
	})
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		Win(w, r)
	})
	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		DeleteCodeCookies(w, r)
		RoomStart(w, r, room)
	})
	http.HandleFunc("/sendData", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		room.broadcastMessage("deaf_" + code)
		http.Redirect(w, r, "/deaftest", http.StatusFound)
	})
	http.HandleFunc("/redirectAll", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		room.broadcastMessage("scattegorie" + code)
		http.Redirect(w, r, "/win", http.StatusFound)
	})
	http.HandleFunc("/sendDataSC", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "BDD.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userid := GetCoockie(w, r, "userId")
		code := GetCoockieCode(w, r, "code")
		roomID, _ := GetRoomByName(db, code)
		room.nbPlayer = checkNbPlayer(db, roomID)

		playersSC++
		log.Println(room.nbPlayer)
		if room.nbPlayer == 1 {
			if playersSC == 1 {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					room.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				playersSC = 0
			}
		} else {
			if playersSC == room.nbPlayer {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					room.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				playersSC = 0
			}
		}
		ScattegoriesRound(w, r, room)
	})
	http.HandleFunc("/deaftest", func(w http.ResponseWriter, r *http.Request) {
		players++
		if Deaftest.nbPlayer == 1 {
			if players == 1 {
				Deaftest.currentPlay++
				players = 0
			}
		} else {
			if players == 1 {
				Deaftest.currentPlay++
			} else if players == Deaftest.nbPlayer {
				players = 0
			}
		}
		DeafTest(w, r, Deaftest)
	})
	http.HandleFunc("/deaftestround", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "BDD.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userid := GetCoockie(w, r, "userId")
		code := GetCoockieCode(w, r, "code")
		roomID, _ := GetRoomByName(db, code)
		if userid == GetCrteatedPlayer(db, roomID) {
			Deaftest.currentMusic = PlaylistConnect(Deaftest.gender)
		}
		last++
		if Deaftest.nbPlayer == 1 {
			if last == 1 {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					Deaftest.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				last = 0
			}
		} else {
			if last == Deaftest.nbPlayer {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					Deaftest.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				last = 0
			}
		}
		DeafTestRound(w, r, Deaftest)
	})
	http.HandleFunc("/deaftestChecker", func(w http.ResponseWriter, r *http.Request) {
		DeafForm(w, r, Deaftest)
	})
	http.HandleFunc("/startPlaying", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		db, err := sql.Open("sqlite3", "BDD.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		roomId, err := GetRoomByName(db, code)
		nbrsUsers := checkNbPlayer(db, roomId)
		if room.game == "scattegories" {
			if err := r.ParseForm(); err != nil {
				log.Fatal(err)
			}
			room.Time = r.FormValue("responseTime")
			room.Letter = selectRandomLetter()
			room.broadcastMessage("data_" + code)
			http.Redirect(w, r, "/scattegories", http.StatusFound)
		} else if room.game == "blindTest" {
			if err := r.ParseForm(); err != nil {
				log.Fatal(err)
			}
			Blindtest.Time = r.FormValue("responseTime")
			Blindtest.gender, Blindtest.nbSong = WaitingFormBT(w, r)
			Blindtest.currentPlay = 0
			Blindtest.nbPlayer = nbrsUsers
			Blindtest.currentBtest = BlindtestManager(Blindtest.gender)
			room.broadcastMessage("blind_" + code)
		} else if room.game == "deafTest" {
			if err := r.ParseForm(); err != nil {
				log.Fatal(err)
			}
			Deaftest.Time = r.FormValue("responseTime")
			Deaftest.gender, Deaftest.nbSong = WaitingForm(w, r)
			Deaftest.currentPlay = 0
			Deaftest.nbPlayer = nbrsUsers
			Deaftest.currentMusic = PlaylistConnect(Deaftest.gender)
			room.broadcastMessage("deaf_" + code)
		}
	})
	http.HandleFunc("/waitingInvit", func(w http.ResponseWriter, r *http.Request) {
		WaitingInvit(w, r, room)
	})
	http.HandleFunc("/sendDataBT", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		room.broadcastMessage("blind_" + code)
		http.Redirect(w, r, "/blindtest", http.StatusFound)
	})
	http.HandleFunc("/blindtest", func(w http.ResponseWriter, r *http.Request) {
		playersBT++
		if Blindtest.nbPlayer == 1 {
			if playersBT == 1 {
				Blindtest.currentPlay++
				playersBT = 0
			}
		} else {
			if playersBT == 1 {
				Blindtest.currentPlay++
			} else if playersBT == Blindtest.nbPlayer {
				playersBT = 0
			}
		}
		BlindTest(w, r, Blindtest)
	})
	http.HandleFunc("/blindtestverif", func(w http.ResponseWriter, r *http.Request) {
		BlindForm(w, r, Blindtest)
	})
	http.HandleFunc("/blindtestround", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "BDD.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userid := GetCoockie(w, r, "userId")
		code := GetCoockieCode(w, r, "code")
		roomID, _ := GetRoomByName(db, code)
		if userid == GetCrteatedPlayer(db, roomID) {
			Blindtest.currentBtest = BlindtestManager(Blindtest.gender)
		}
		lastBT++
		if Blindtest.nbPlayer == 1 {
			if lastBT == 1 {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					Blindtest.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				lastBT = 0
			}
		} else {
			if lastBT == Blindtest.nbPlayer {
				roomCreator := GetCrteatedPlayer(db, roomID)
				if roomCreator == userid {
					Blindtest.finish = true
				}
				room.broadcastMessage(strconv.Itoa(roomCreator))
				lastBT = 0
			}
		}
		BlindTestRound(w, r, Blindtest)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(room, w, r)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}

func HandleWebSocket(room *Room, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	room.clients[ws] = true
	room.register <- ws
	defer func() { room.unregister <- ws }()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
	}
}
