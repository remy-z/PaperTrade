package main

import "encoding/json"

// constants for routing on both front and backend
// naming is from the client viewpoint
// send_message are messages sent by client, recieved by server
// recieve_message are messages recieved by client, sent by server
const (
	EventSendSub      = "send_sub"
	EventSendUnsub    = "send_unsub"
	EventRecievePrice = "recieve_price"
	EventRecieveError = "recieve_err"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendSubEvent struct {
	Symbol string `json:"symbol"`
}

type RecieveError struct {
	Message string `json:"message"`
}

type SendUnsubEvent struct {
	Symbol string `json:"symbol"`
}

type ReceivePriceEvent struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type EventHandler func(event Event, c *Client) error
