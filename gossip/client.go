package gossip

import (
	"log"
	"net/http"
	"time"
)

type Client struct {
	nodes    []nodeDescriptor
	maxNodes int

	messages    []messageDescriptor
	maxMessages int

	shutdown chan bool

	//port uint64

	aliveTimeout   time.Duration
	messageTimeout time.Duration
}

// constructor
func New(maxNodes int, maxMessages int, aliveTimeout time.Duration, messageTimeout time.Duration) *Client {
	return &Client{
		nodes:          make([]nodeDescriptor, maxNodes),
		maxNodes:       maxNodes,
		messages:       make([]messageDescriptor, maxNodes),
		maxMessages:    maxMessages,
		shutdown:       make(chan bool),
		aliveTimeout:   aliveTimeout,
		messageTimeout: messageTimeout,
	}
}

func (gc *Client) messageHandler(writer http.ResponseWriter, request *http.Request) {
	// accept and store incoming messages
	// if queue is full, then delete the oldest one
}

func (gc *Client) aliveHandler(writer http.ResponseWriter, request *http.Request) {
	// accept and process incoming keepalives
	// if queue is full, then delete the oldest node descriptor
}

func (gc *Client) sendMessages() {
	// 1. select a random message and node descriptor
	// 2. send message to node, and request a message from node
	// 3. merge reply into own message slice
}

func (gc *Client) sendAlive() {
	// select a random node descriptor and send keepalive
}

func (gc *Client) mergeMessage(message string) {
	// utility method
}

func (gc *Client) mergeNode(descriptor nodeDescriptor) {
	// utility method
}

func (gc *Client) Send(message string) error {
	// send a new message to the network

	return nil
}

func (gc *Client) sendLoop() {
	aliveTicker := time.NewTicker(gc.aliveTimeout)
	messageTicker := time.NewTicker(gc.messageTimeout)

loop:
	for {
		select {
		case <-aliveTicker.C:
			// send keepalive
		case <-messageTicker.C:
			// send message
		case <-gc.shutdown:
			break loop
		}
	}
	aliveTicker.Stop()
	messageTicker.Stop()
	log.Println("send loop shut down")
}

func (gc *Client) Shutdown() {
	close(gc.shutdown)
}

func (gc *Client) Run() error {
	//http.HandleFunc("/alive/", gc.aliveHandler)
	//http.HandleFunc("/message/", gc.messageHandler)

	// use these with "encoding/gob"
	go gc.recvMessages()
	go gc.recvAlives()

	go gc.sendMessages()
	go gc.sendAlives()
	//return http.ListenAndServe(":8080", nil)
}
