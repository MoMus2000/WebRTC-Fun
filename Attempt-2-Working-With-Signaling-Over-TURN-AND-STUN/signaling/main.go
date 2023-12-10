package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	PORT = ":3000"
)

var (
	FirstClient    *websocket.Conn
	SecondClient   *websocket.Conn
	EncryptedOffer string
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mutex = sync.Mutex{}
)

type Payload struct {
	Action  string
	Message string
	Type    string `json:"type"`
}

type ClientPayload struct {
	Type string `json:"type"`
}

type Offer struct {
	Offer string
}

type IceCandidate struct {
	IceCandidate string
}

type Answer struct {
	Answer string
}

type Connection struct {
	Websocket     *websocket.Conn
	Offer         Offer
	IceCandidates []IceCandidate
	Answer        Answer
}

var ConnectionMapper map[string][]*Connection

func checkIfMapIsFull() {
	for {
		if len(ConnectionMapper["DefaultRoom"]) == 2 {
			time.Sleep(10 * time.Second)
			fmt.Println("Waiting to send ...")
			payload := Payload{
				Action:  "data",
				Message: ConnectionMapper["DefaultRoom"][0].Offer.Offer,
				Type:    "offer",
			}
			b, _ := json.Marshal(payload)
			fmt.Println(string(b))
			fmt.Println("BEFORE WRITING ... ", ConnectionMapper["DefaultRoom"][1].Answer)
			ConnectionMapper["DefaultRoom"][1].Websocket.WriteMessage(websocket.TextMessage, b)
			time.Sleep(5 * time.Second)
			fmt.Println("AFTER WRITING ... ", ConnectionMapper["DefaultRoom"][1].Answer)
			for {
				// if err != nil {
				// }
				// fmt.Println("ANSWER", string(answer))

				payload := Payload{
					Action:  "data",
					Message: ConnectionMapper["DefaultRoom"][1].Answer.Answer,
					Type:    "answer",
				}
				b, _ := json.Marshal(payload)
				ConnectionMapper["DefaultRoom"][0].Websocket.WriteMessage(websocket.TextMessage, b)
				break
			}

			for _, iceCandidate := range ConnectionMapper["DefaultRoom"][0].IceCandidates {
				payload := Payload{
					Action:  "data",
					Message: iceCandidate.IceCandidate,
					Type:    "candidate",
				}
				b, _ := json.Marshal(payload)
				fmt.Println(string(b))
				ConnectionMapper["DefaultRoom"][1].Websocket.WriteMessage(websocket.TextMessage, b)
			}
			for _, iceCandidate := range ConnectionMapper["DefaultRoom"][1].IceCandidates {
				payload := Payload{
					Action:  "data",
					Message: iceCandidate.IceCandidate,
					Type:    "candidate",
				}
				b, _ := json.Marshal(payload)
				fmt.Println(string(b))
				ConnectionMapper["DefaultRoom"][0].Websocket.WriteMessage(websocket.TextMessage, b)
			}
			break
		}
	}
}

func GuideMessages(conn *websocket.Conn) {
	offer := Offer{}
	randomNumber := rand.Intn(100)
	fmt.Println(randomNumber)
	iceCandidates := []IceCandidate{}
	connection := Connection{}
	connection.Websocket = conn
	answer := Answer{}

	ConnectionMapper["DefaultRoom"] = append(ConnectionMapper["DefaultRoom"], &connection)

	fmt.Println("Number of Connections ", len(ConnectionMapper["DefaultRoom"]))

	for {
		_, Message, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println(err)
		}
		Type := ClientPayload{}
		json.Unmarshal(Message, &Type)
		switch Type.Type {
		case "offer":
			offer.Offer = string(Message)
			connection.Offer = offer
		case "candidate":
			iceCandidates = append(iceCandidates, IceCandidate{string(Message)})
			connection.IceCandidates = append(connection.IceCandidates, IceCandidate{string(Message)})
		case "answer":
			fmt.Println("ANSSSWEERR")
			fmt.Println(string(Message))
			answer.Answer = string(Message)
			connection.Answer = answer
		}
	}
}

var i int

func Gateway(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	if i == 0 {
		data := Payload{Action: "ready", Message: ""}
		bytes, _ := json.Marshal(data)

		conn.WriteMessage(websocket.TextMessage, bytes)

		i += 1
	}

	GuideMessages(conn)
}

func main() {
	ConnectionMapper = make(map[string][]*Connection)
	go checkIfMapIsFull()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ws", Gateway)

	http.ListenAndServe(":9999", r)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
