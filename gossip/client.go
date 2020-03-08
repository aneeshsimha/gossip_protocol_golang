package gossip

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

const (
	CHANSIZE = 5
)

type Client struct {
	nodes    []nodeDescriptor
	maxNodes int

	messages    []messageDescriptor
	maxMessages int

	shutdown    chan bool
	aliveChan   chan nodeDescriptor
	messageChan chan messageDescriptor

	alivePort   string
	messagePort string

	aliveTimeout   time.Duration
	messageTimeout time.Duration
}

// constructor
func New(maxNodes int, maxMessages int, alivePort string, messagePort string, aliveTimeout time.Duration, messageTimeout time.Duration) *Client {
	return &Client{
		nodes:          make([]nodeDescriptor, maxNodes),
		maxNodes:       maxNodes,
		messages:       make([]messageDescriptor, maxNodes),
		maxMessages:    maxMessages,
		shutdown:       make(chan bool),
		aliveChan:      make(chan nodeDescriptor, CHANSIZE),
		messageChan:    make(chan messageDescriptor, CHANSIZE),
		alivePort:      alivePort,
		messagePort:    messagePort,
		aliveTimeout:   aliveTimeout,
		messageTimeout: messageTimeout,
	}
}

//func (gc *Client) messageHandler(writer http.ResponseWriter, request *http.Request) {
//	// accept and store incoming messages
//	// if queue is full, then delete the oldest one
//}
//
//func (gc *Client) aliveHandler(writer http.ResponseWriter, request *http.Request) {
//	// accept and process incoming keepalives
//	// if queue is full, then delete the oldest node descriptor
//}

func (gc *Client) sendMessages() {
	// 1. select a random message and node descriptor
	// 2. send message to node, and request a message from node
	// 3. merge reply into own message slice
}

func (gc *Client) recvMessages() {
	// TODO
}

func (gc *Client) handleMessage(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	msg := StringPayload{}
	dec.Decode(&msg)

	gc.messageChan <- msg.Message
}

func (gc *Client) sendAlives() {
	// select a random node descriptor and send keepalive

	// add own ip + current time to the copy of the node descriptor list before sending
}

func (gc *Client) recvAlives() {
	listener, _ := net.Listen("tcp", gc.alivePort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go gc.handleAlive(conn)
	}
}

func (gc *Client) handleAlive(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	kap := KeepAlivePayload{}
	dec.Decode(&kap)

	for _, desc := range kap.KnownNodes {
		// do stuff
		gc.aliveChan <- desc
	}
}

func (gc *Client) mergeNode(descriptor nodeDescriptor) {
	// utility method
	// TODO
}

func (gc *Client) mergeMessage(message messageDescriptor) {
	// utility method
	// TODO
}

func (gc *Client) process() {
	// process messages that have been sent down the various channels
	for {
		select {
		case <-gc.shutdown:
			return
		case node := <-gc.aliveChan:
			// merge the node in

			gc.mergeNode(node) // TODO
		case message := <-gc.messageChan:
			// merge the message in

			gc.mergeMessage(message) // TODO
		default:
			// do nothing lol
		}
	}
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

// exposed methods:

func (gc *Client) Send(message string) error {
	// send a new message to the network
	// TODO

	return nil
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

	go gc.process()
	//return http.ListenAndServe(":8080", nil)
}
