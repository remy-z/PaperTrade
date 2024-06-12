package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	sync.RWMutex
	clients   ClientList
	symbolMap SymbolMap
	handlers  map[string]EventHandler
	prices    map[string]float32
	fhClient  *FhClient
}

type ClientList map[*Client]bool
type SymbolMap map[string]ClientList

func NewManager(fhClient *FhClient) *Manager {
	// Create an FhClient instance
	m := &Manager{
		clients:   make(ClientList),
		symbolMap: make(SymbolMap),
		handlers:  make(map[string]EventHandler),
		prices:    make(map[string]float32),
		fhClient:  fhClient,
	}
	m.setupEventhandlers()
	fhClient.setupManager(m)
	go fhClient.readMessages()
	go fhClient.writeMessages()
	return m
}

// http request handler for /ws, upgrades connection to websocket
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	fmt.Println(conn)
	if err != nil {
		fmt.Println(err)
	}

	client := NewClient(conn, m)
	m.addClient(client)

	// Start read and write go routines
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) setupEventhandlers() {
	m.handlers[EventSendSub] = SendSubHandler
	m.handlers[EventSendUnsub] = SendUnsubHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("event type not found")
	}
}

func SendSubHandler(event Event, c *Client) error {
	var sendSubEvent SendSubEvent
	err := json.Unmarshal(event.Payload, &sendSubEvent)
	if err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	err = c.manager.subToSymbol(c, sendSubEvent.Symbol)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}
	return nil
}

// check symbol map if symbol already subbed
// sub to symbol in fh if needed
// add client to clientlist in symbol map
func (m *Manager) subToSymbol(c *Client, symbol string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.symbolMap[symbol]; !ok {
		m.fhClient.subscribeTicker(symbol)
		m.symbolMap[symbol] = make(ClientList)
	}
	m.symbolMap[symbol][c] = true
	return nil
}

// remove client for clientlist in symbolmap
// if clientlist empty, delete key and unsub in fh
func SendUnsubHandler(event Event, c *Client) error {
	var sendUnsubEvent SendUnsubEvent
	err := json.Unmarshal(event.Payload, &sendUnsubEvent)
	if err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	err = c.manager.unsubToSymbol(c, sendUnsubEvent.Symbol)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}
	return nil
}

func (m *Manager) unsubToSymbol(c *Client, symbol string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.symbolMap[symbol], c)
	if len(m.symbolMap[symbol]) == 0 {
		m.fhClient.unsubscribeTicker(symbol)
		delete(m.symbolMap, symbol)
	}
	return nil
}

// add client to manager's ClientList
// add client to chatroom map
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

// close the client's connection and remove client from memory completely
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	delete(m.clients, client)
	client.conn.Close()
}

// parse the object that finnhub sends us, broadcast the prices of assets to clients that are subbed
func (m *Manager) parseFhPrice(msg SendFhTrade) error {
	symbolToPrice := make(map[string]float64)
	for _, trade := range msg.Data {
		symbolToPrice[trade.S] = trade.P
	}
	for s, p := range symbolToPrice {
		priceEvent := ReceivePriceEvent{
			Symbol: s,
			Price:  p,
		}
		m.broadcastPrice(priceEvent)
	}
	return nil
}
func (m *Manager) broadcastPrice(priceEvent ReceivePriceEvent) error {
	data, err := json.Marshal(priceEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventRecievePrice,
	}

	//write messages to egress channel on clients
	for client := range m.symbolMap[priceEvent.Symbol] {
		client.egress <- outgoingEvent
	}
	return nil
}
