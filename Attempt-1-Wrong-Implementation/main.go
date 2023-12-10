package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	PORT = ":3000"
)

var (
	OfferClient    *websocket.Conn
	AnswerClient   *websocket.Conn
	EncryptedOffer string
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func OfferWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	OfferClient = conn
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Infinite loop to read messages from the WebSocket connection
	i := 0
	for {
		// Read message from the WebSocket connection
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		if len(string(p)) > 0 && i == 0 {
			EncryptedOffer = string(p)
			fmt.Println(string(EncryptedOffer))
			AnswerClient.WriteMessage(websocket.TextMessage, []byte(EncryptedOffer))
			i += 1
		} else if len(string(p)) > 0 && i == 1 {
			Answer := string(p)
			OfferClient.WriteMessage(websocket.TextMessage, []byte(Answer))
			i += 1
		}
		// Print the received message
		// fmt.Printf("Received message: %s\n", p)

		// Echo the message back to the client
		// err = conn.WriteMessage(messageType, p)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
}

func AnswerWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	AnswerClient = conn
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Infinite loop to read messages from the WebSocket connection
	for {
		// Read message from the WebSocket connection
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		if len(string(p)) > 0 {
			OfferClient.WriteMessage(websocket.TextMessage, p)
		}

		// Print the received message
		fmt.Printf("Received message: %s\n", p)

		// // Echo the message back to the client
		// err = conn.WriteMessage(messageType, p)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
}

func ServeOfferJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/js/offer.js")
}

func ServeAnswerJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/js/answer.js")
}

func ServeOfferPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/offer.html")
}

func ServeAnswerPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/answer.html")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", ServeOfferPage)
	r.HandleFunc("/js/offer", ServeOfferJs)
	r.HandleFunc("/js/answer", ServeAnswerJs)
	r.HandleFunc("/answer", ServeAnswerPage)
	r.HandleFunc("/offer/ws", OfferWebSocket)
	r.HandleFunc("/answer/ws", AnswerWebSocket)

	fmt.Println("Listening on PORT ", PORT)
	http.ListenAndServe(PORT, r)
}
