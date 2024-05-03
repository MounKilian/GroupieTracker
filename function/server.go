package groupieTracker

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var questions = []Question{}

var state = false

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

	http.HandleFunc("/", Home)
	http.HandleFunc("/checkUser", func(w http.ResponseWriter, r *http.Request) {
		Formulaire(w, r)
	})
	http.HandleFunc("/scattegories", func(w http.ResponseWriter, r *http.Request) {
		room.broadcastMessage("data_" + room.letter)
		Scattegories(w, r, room.letter)
	})
	http.HandleFunc("/scattegoriesChecker", func(w http.ResponseWriter, r *http.Request) {
		userid := GetCoockie(w, r, "userId")
		if userid == 3 {
			if !state {
				room.broadcastMessage("end")
				state = true
				response := ScattegoriesForm(w, r)
				questions = append(questions, response)
				fmt.Println(questions)
			} else {
				ScattegoriesForm(w, r)
			}
		} else {
			response := ScattegoriesForm(w, r)
			questions = append(questions, response)
			fmt.Println(questions)
		}
	})
	http.HandleFunc("/verification", func(w http.ResponseWriter, r *http.Request) {
		ScattegoriesVerification(w, r, questions[0])
		ScattegoriesVerification(w, r, questions[1])
	})
	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		userid := GetCoockie(w, r, "userId")
		if userid == 3 {
			room.letter = selectRandomLetter()
			room.broadcastMessage(room.letter)
		}
		Waiting(w, r)
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
