package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sudonite/TollCalculator/types"
)

const sendDataNum = 20
const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second * 5

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genOBUID() int {
	return rand.Intn(math.MaxInt)
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < sendDataNum; i++ {
			obuID := genOBUID()
			lat, long := genLatLong()

			data := types.OBUData{
				OBUID: obuID,
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}
