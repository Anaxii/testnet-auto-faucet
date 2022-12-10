package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testnet-autofaucet/util"
)

var SocketChannels =  make(map[string]chan interface{})

func Log(data interface{}) {
	for  k := range SocketChannels {
		SocketChannels[k] <- data
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 30,
	WriteBufferSize: 1024 * 30,
}

func getWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	data := make(chan interface{})
	id := util.RandStringRunes(20)
	SocketChannels[id] = data
	reader(ws, data, id)
}

func reader(conn *websocket.Conn, dataChannel chan interface{}, id string) {
	go handleIncomingMessages(conn, id)
	for {
		select {
		case d := <-dataChannel:
			response := map[string]interface{}{"status": "log", "data": d}
			data, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				delete(SocketChannels, id)
				log.Println(err)
				return
			}
		}
	}
}

func handleIncomingMessages(conn *websocket.Conn, id string) {
	x := 0
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(SocketChannels, id)
			return
		}
		response := map[string]string{"status": "Connection to Puffin Auto-Faucet established"}
		data, _ := json.Marshal(response)
		if x == 0 {
			x++
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				delete(SocketChannels, id)
				return
			}
		}

		if string(msg) == "ping" {
			response = map[string]string{"status": "pong"}
			data, _ = json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				delete(SocketChannels, id)
				return
			}
		}

	}
}