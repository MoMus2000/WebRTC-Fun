package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mutex         = sync.Mutex{}
	AnswerChannel chan string
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
	Value string
}

type IceCandidate struct {
	Value string
}

type Connection struct {
	Websocket     *websocket.Conn
	Offer         Offer
	IceCandidates []IceCandidate
}

var ConnectionMapper map[string][]*Connection

func checkIfRoomMapIsFull(key string) {
	for {
		fmt.Println(key)
		mutex.Lock()
		length := len(ConnectionMapper[key])
		mutex.Unlock()
		if length == 2 {
			payload := Payload{
				Action:  "data",
				Message: ConnectionMapper[key][0].Offer.Value,
				Type:    "offer",
			}
			b, _ := json.Marshal(payload)
			ConnectionMapper[key][1].Websocket.WriteMessage(websocket.TextMessage, b)
			var ans string
			ans = <-AnswerChannel

			if ans != "" {
				payload := Payload{
					Action:  "data",
					Message: ans,
					Type:    "answer",
				}
				b, _ := json.Marshal(payload)
				ConnectionMapper[key][0].Websocket.WriteMessage(websocket.TextMessage, b)
			}

			// Observation found that a small delay is needed on mobile after sending the webrtc answer
			time.Sleep(1 * time.Second)

			mutex.Lock()
			for _, iceCandidate := range ConnectionMapper[key][0].IceCandidates {
				payload := Payload{
					Action:  "data",
					Message: iceCandidate.Value,
					Type:    "candidate",
				}
				b, _ := json.Marshal(payload)
				ConnectionMapper[key][1].Websocket.WriteMessage(websocket.TextMessage, b)
			}
			mutex.Unlock()

			mutex.Lock()
			for _, iceCandidate := range ConnectionMapper[key][1].IceCandidates {
				payload := Payload{
					Action:  "data",
					Message: iceCandidate.Value,
					Type:    "candidate",
				}
				b, _ := json.Marshal(payload)
				ConnectionMapper[key][0].Websocket.WriteMessage(websocket.TextMessage, b)
			}
			mutex.Unlock()
			break
		}
	}
}

func GuideMessages(conn *websocket.Conn, roomKey string) {
	offer := Offer{}
	iceCandidates := []IceCandidate{}
	connection := Connection{}
	connection.Websocket = conn

	mutex.Lock()
	ConnectionMapper[roomKey] = append(ConnectionMapper[roomKey], &connection)
	mutex.Unlock()

	fmt.Println("Number of Connections ", len(ConnectionMapper[roomKey]))

	for {
		_, Message, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println(err)
		}
		Type := ClientPayload{}
		json.Unmarshal(Message, &Type)
		switch Type.Type {
		case "offer":
			offer.Value = string(Message)
			mutex.Lock()
			connection.Offer = offer
			mutex.Unlock()
		case "candidate":
			iceCandidates = append(iceCandidates, IceCandidate{string(Message)})
			mutex.Lock()
			connection.IceCandidates = append(connection.IceCandidates, IceCandidate{string(Message)})
			mutex.Unlock()
		case "answer":
			mutex.Lock()
			AnswerChannel <- string(Message)
			mutex.Unlock()
		}
	}
}

func Gateway(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var roomKey string = ""
	urlParams := r.URL.Query()

	key, ok := urlParams["roomId"]
	fmt.Println(urlParams)

	fmt.Println(key[0])

	if ok {
		roomKey = key[0]
	} else {
		roomKey = "DefaultRoom"
	}

	if len(ConnectionMapper[roomKey]) == 0 {
		data := Payload{Action: "ready", Message: ""}
		bytes, _ := json.Marshal(data)
		conn.WriteMessage(websocket.TextMessage, bytes)
		go checkIfRoomMapIsFull(roomKey)
	}

	GuideMessages(conn, roomKey)
}

func main() {
	AnswerChannel = make(chan string)
	ConnectionMapper = make(map[string][]*Connection)
	ConnectionMapper["DefaultRoom"] = make([]*Connection, 0)

	r := mux.NewRouter().StrictSlash(true)

	// provide room ID, otherwise its the default room
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
