package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sudonite/TollCalculator/types"
)

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

func genOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	obuIDS := genOBUIDS(20)

	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()

			data := types.OBUData{
				OBUID: obuIDS[i],
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
