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

type Deaftest struct {
	gender       string
	nbSong       int
	currentMusic Music
	currentPlay  int
	nbPlayer     int
}

type Room struct {
	id         string
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
	letter     string
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
	var Deaftest Deaftest
	questionsMap = make(map[string][]Question)

	http.HandleFunc("/", Home)
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		Formulaire(w, r)
	})
	http.HandleFunc("/landingPage", func(w http.ResponseWriter, r *http.Request) {
		LandingPage(w, r)
	})
	http.HandleFunc("/scattegories", func(w http.ResponseWriter, r *http.Request) {
		Scattegories(w, r, room.letter)
	})
	http.HandleFunc("/scattegoriesChecker", func(w http.ResponseWriter, r *http.Request) {
		buttonValue := r.FormValue("button-value")
		code := GetCoockieCode(w, r, "code")
		userid := GetCoockie(w, r, "userId")
		questions := questionsMap[code]
		log.Println("hello")
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
		log.Println(questions)
		ScattegoriesVerification(w, r, questions)
	})
	http.HandleFunc("/verificationChecker", func(w http.ResponseWriter, r *http.Request) {
		ScattegoriesVerificationChecker(w, r)
	})
	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		Waiting(w, r, room)
	})
	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		RoomStart(w, r)
	})
	http.HandleFunc("/sendData", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		room.broadcastMessage("deaf_" + code)
	})
	http.HandleFunc("/deaftest", func(w http.ResponseWriter, r *http.Request) {
		players++
		if players == 1 {
			Deaftest.currentPlay++
		} else if players == 2 {
			players = 0
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
		DeafTestRound(w, r, Deaftest)
	})
	http.HandleFunc("/deaftestChecker", func(w http.ResponseWriter, r *http.Request) {
		DeafForm(w, r, Deaftest)
	})
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		Win(w, r, Deaftest.currentMusic)
	})
	http.HandleFunc("/startPlaying", func(w http.ResponseWriter, r *http.Request) {
		code := GetCoockieCode(w, r, "code")
		if game == "scattegories" {
			room.letter = selectRandomLetter()
			room.broadcastMessage("data_" + code)
			http.Redirect(w, r, "/scattegories", http.StatusFound)
		} else if game == "blindTest" {
			http.Redirect(w, r, "/room", http.StatusFound)
		} else if game == "deafTest" {
			Deaftest.gender, Deaftest.nbSong = WaitingForm(w, r)
			Deaftest.currentMusic = PlaylistConnect(Deaftest.gender)
			// Deaftest.currentPlay++
			room.broadcastMessage("deaf_" + code)
		}
	})
	http.HandleFunc("/waitingInvit", func(w http.ResponseWriter, r *http.Request) {
		WaitingInvit(w, r, room)
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
