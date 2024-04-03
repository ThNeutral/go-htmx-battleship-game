package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	empty  = 0
	placed = 1
	missed = 2
	hit    = 3
)

type Fields struct {
	FirstField  [100]int `json:"first_field"`
	SecondField [100]int `json:"second_field"`
}

type IncomingMessage struct {
	From             int   `json:"from"`
	ChangeTo         int   `json:"change_to"`
	ChangedPositions []int `json:"changed_positions"`
}

var currentTurn = 1
var connections = make(map[*websocket.Conn]int)
var fields Fields
var channel = make(chan *IncomingMessage)

func getNextPlayerID() int {
	if len(connections) == 0 {
		return 1
	}
	count := 0
	for _, value := range connections {
		count += value
	}
	if count == 1 {
		return 2
	} else {
		return 1
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	if len(connections) >= 2 {
		log.Println("Lobby is already full")
		return
	}

	log.Println("Encountered attempt to create WebSocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() { conn.Close(); delete(connections, conn) }()

	connections[conn] = getNextPlayerID()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error while reading message: ", err)
			return
		}

		var incMess *IncomingMessage
		err = json.Unmarshal(message, incMess)
		if err != nil {
			log.Println("Failed to unmarshal: ", err.Error())
			return
		}

		if incMess.From != currentTurn {
			continue
		}

		channel <- incMess

		log.Println("Recieved message: ", string(message))
	}
}

func handleIncomingMessages() {
	for {
		incMess := <-channel
		if incMess.From == 1 {
			for _, val := range incMess.ChangedPositions {
				if incMess.ChangeTo == placed {
					fields.FirstField[val] = placed
				} else if incMess.ChangeTo == hit {
					if fields.SecondField[val] == placed {
						fields.SecondField[val] = hit
					} else {
						fields.SecondField[val] = missed
					}
				}
			}
		} else {
			for _, val := range incMess.ChangedPositions {
				if incMess.ChangeTo == placed {
					fields.SecondField[val] = placed
				} else if incMess.ChangeTo == hit {
					if fields.FirstField[val] == placed {
						fields.FirstField[val] = hit
					} else {
						fields.FirstField[val] = missed
					}
				}
			}
		}
		var f1, f2 Fields

		for conn, playerID := range connections {
			if playerID == 1 {
				f1.FirstField = fields.FirstField
				for i, val := range fields.SecondField {
					if val == placed {
						f1.SecondField[i] = empty
					} else {
						f1.SecondField[i] = val
					}
				}
				json, err := json.Marshal(f1)
				if err != nil {
					log.Println("Failed to marshal: ", err.Error())
				}
				conn.WriteMessage(websocket.TextMessage, json)
			} else {
				f2.SecondField = fields.SecondField
				for i, val := range fields.FirstField {
					if val == placed {
						f1.FirstField[i] = empty
					} else {
						f1.FirstField[i] = val
					}
				}
				json, err := json.Marshal(f1)
				if err != nil {
					log.Println("Failed to marshal: ", err.Error())
				}
				conn.WriteMessage(websocket.TextMessage, json)
			}
		}
	}
}

func main() {
	http.HandleFunc("/connect", handleConnect)

	go handleIncomingMessages()

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
