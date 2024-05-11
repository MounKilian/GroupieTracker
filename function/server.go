package groupieTracker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Blindtest struct {
	gender       string
	nbSong       int
	currentBtest Btest
	currentPlay  int
}

func Server() {
	room := NewRoom()
	go room.Start()
	var currentUser User
	var Blindtest Blindtest

	http.HandleFunc("/", Home)
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		currentUser = Formulaire(w, r)
		fmt.Println(currentUser)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(room, w, r)
	})

	http.HandleFunc("/blindtest", func(w http.ResponseWriter, r *http.Request) {
		Blindtest.gender, Blindtest.nbSong = "6EyNHMMJtumWbxWpWN5AJW", 5
		Blindtest.currentBtest = BlindtestManager(Blindtest.gender)
		BlindTest(w, r, Blindtest)
	})

	http.HandleFunc("/blindtestverif", func(w http.ResponseWriter, r *http.Request) {
		BlindForm(w, r, Blindtest)
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
