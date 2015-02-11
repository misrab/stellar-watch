package main 

import (
	"log"
	// "bytes"
	"strconv"

	"encoding/json"

	"github.com/misrab/web"
)


const (
	url = "wss://live.stellar.org:9001"
	body = `{ "command": "subscribe", "streams": [ "transactions" ] }`
)

func main() {
	log.Println("Running main...")


	dbmap := SetupDB()
	amount_chan := make(chan int)
	go handleAmounts(amount_chan, dbmap)


	go listenToRequests(dbmap)


	// listen to websocket
	response_chan := make(chan []byte)
	go web.ListenToSocket(url, body, response_chan)
	response := make(map[string]interface{})
	// transaction := make(map[string]interface{})
	for data := range response_chan {
		// get amount
		json.Unmarshal(data, &response)
		transaction := response["transaction"]
		if transaction == nil { 
			log.Println("Transaction nil, skipping...")
			continue 
		}
		amount_raw := transaction.(map[string]interface{})["Amount"]
		if amount_raw == nil { 
			log.Println("Amount nil, skipping...")
			continue 
		}
		amount , err := strconv.Atoi(amount_raw.(string))
		if err != nil { 
			log.Println("Amount not integer, skipping...")
			continue 
		}

		amount_chan <- amount
	}
}