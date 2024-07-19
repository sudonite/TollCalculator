package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sudonite/TollCalculator/types"
)

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU Connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		fmt.Printf("Received OBU Data: [%d] :: <lat %.2f, long %.2f>\n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data
	}
}

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}
