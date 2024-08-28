package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type FhClient struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan Event
}

type SendFhTrade struct {
	Data []struct {
		P float64 `json:"p"` //price
		S string  `json:"s"` //symbol
		T int64   `json:"t"` //time
		V float64 `json:"v"` //volume
	} `json:"data"`
	Type string `json:"type"`
}

func NewFinnHubClient(token string) *FhClient {
	conn, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token="+token, nil)
	if err != nil {
		panic(err)
	}
	return &FhClient{
		conn: conn,
	}
}

func (f *FhClient) setupManager(manager *Manager) error {
	f.manager = manager
	return nil
}

func (f *FhClient) readMessages() error {
	defer f.conn.Close()

	var msg SendFhTrade
	for {
		err := f.conn.ReadJSON(&msg)
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
			break
		}
		if msg.Type == "trade" {
			f.manager.parseFhPrice(msg)
		} else {
			fmt.Print(msg)
		}
	}
	return nil
}

func (f *FhClient) writeMessages() {
	for {
		//write message to websocket connection whenever a messages comes from egress channel
		message, ok := <-f.egress
		if !ok {
			//Server write message to client when server has egress issues that connection has to be closed
			if err := f.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
				fmt.Println("connection closed: ", err)
			}
			return
		}

		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := f.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Printf("failed to send message: %v", err)
		}
		fmt.Println("message sent")
	}
}

func (f *FhClient) subscribeTicker(s string) error {
	f.writeToFh(map[string]string{"type": "subscribe", "symbol": s})
	return nil
}

func (f *FhClient) unsubscribeTicker(s string) error {
	f.writeToFh(map[string]string{"type": "unsubscribe", "symbol": s})
	return nil
}

func (f *FhClient) writeToFh(message map[string]string) error {
	msg, _ := json.Marshal(message)
	f.conn.WriteMessage(websocket.TextMessage, msg)
	return nil
}
