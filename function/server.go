package groupieTracker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var questions = []Question{}

func Server() {
	room := NewRoom()
	go room.Start()
	var letter string
	var currentMusic Music

	http.HandleFunc("/", Home)
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		Formulaire(w, r)
	})
	http.HandleFunc("/scattegories", func(w http.ResponseWriter, r *http.Request) {
		letter = selectRandomLetter()
		Scattegories(w, r, letter)
	})
	http.HandleFunc("/scattegoriesChecker", func(w http.ResponseWriter, r *http.Request) {
		response := ScattegoriesForm(w, r)
		questions = append(questions, response)
		fmt.Println(questions)
	})
	http.HandleFunc("/verification", func(w http.ResponseWriter, r *http.Request) {
		ScattegoriesVerification(w, r, questions[0])
	})
	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		Waiting(w, r)
	})
	http.HandleFunc("/deaftest", func(w http.ResponseWriter, r *http.Request) {
		currentMusic = PlaylistConnect()
		DeafTest(w, r, currentMusic)
	})
	http.HandleFunc("/deaftestverif", func(w http.ResponseWriter, r *http.Request) {
		CheckDeafTest(w, r, currentMusic)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(room, w, r)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}

type Room struct {
	id         string
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
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

func HandleWebSocket(room *Room, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

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
