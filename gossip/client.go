package gossip

import (
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"time"
	"bufio"
    "fmt"
    "io/ioutil"
    "os"

	"coen317/gossip/counter"
)

const (
	CHANSIZE = 5
)

type Client struct {
	// collision chance in a 64-bit ID space is n^2 / 2^65
	id uint64

	self     nodeDescriptor
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

	counter *counter.Counter // threadsafe counter type
}

// constructor
func New(maxNodes int, maxMessages int, alivePort string, messagePort string, aliveTimeout time.Duration, messageTimeout time.Duration) *Client {
	id := rand.Uint64()
	for id < 100 {
		id = rand.Uint64() // 0-99 are reserved
	}
	return &Client{
		id:             id,
		self:           nodeDescriptor{},
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
		counter:        counter.New(),
	}
}

func (gc *Client) sendMessages() {
	// 1. select a random message and node Descriptor
	// 2. send message to node, and request a message from node
	// 3. merge reply into own message slice

	// select a random node Descriptor and send random message

	messageTicker := time.NewTicker(gc.messageTimeout)
	defer messageTicker.Stop()
	defer log.Println("send loop shut down")

	// loop forever
	for {
		select {
		case <-messageTicker.C: // do every interval
			//  choose a random known node descriptor
			randomNode := randomNeighbor(gc.nodes)
			// TODO

			//  choose a random stored message
			//  turn the messageDescriptor into a stringPacket
			//  send the stringPacket
		case <-gc.shutdown:
			return
		}
	}
}

func (gc *Client) recvMessages() {
	listener, _ := net.Listen("tcp", gc.messagePort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go gc.handleMessage(conn)
	}
}

func (gc *Client) handleMessage(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	msg := StringPayload{}
	dec.Decode(&msg)

	if msg.Message.ID == gc.id {
		return // don't bother adding own originating messages
	}

	gc.messageChan <- msg.Message
}

// goroutine
func (gc *Client) messageLoop() {
	for {
		desc := <-gc.messageChan
		insertMessage(gc.messages[:], desc)
	}
}

func (gc *Client) sendAlives() {
	// select a random node Descriptor and send keepalive

	aliveTicker := time.NewTicker(gc.aliveTimeout)
	defer aliveTicker.Stop()
	defer log.Println("send loop shut down")

	// loop forever
	for {
		select {
		case <-aliveTicker.C: // do every interval
			gc.sendAlive() // package this into a single function because it has a defer
		case <-gc.shutdown:
			return
		}
	}
}

func (gc *Client) sendAlive() {
	// choose a random known node descriptor
	// add own ip + current time to the copy of the node Descriptor list before sending

	randomNode := randomNeighbor(gc.nodes)
	conn, err := net.Dial("tcp", randomNode.Address.String()+":"+gc.alivePort)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close() // this is why this is its own function

	// make/get self node
	me := newNodeDescriptor(net.ParseIP(conn.LocalAddr().String()), time.Now(), gc.id, <-gc.counter.Count)
	kap := preparePayload(gc.nodes, me)

	enc := gob.NewEncoder(conn)
	_ = enc.Encode(kap)
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
		if desc.ID == gc.id {
			continue // don't bother adding own originating messages
		}
		logFile(f, desc)
		gc.aliveChan <- desc
	}
}

// goroutine
func (gc *Client) aliveLoop() {
	for {
		desc := <-gc.aliveChan
		insertNode(gc.nodes[:], desc)
	}
}

// replaced by various things
//func (gc *Client) mergeNode(descriptor nodeDescriptor) {
//	// utility method
//	// TODO
//	insertNode(gc.nodes[:], descriptor)
//}
//
//func (gc *Client) mergeMessage(message messageDescriptor) {
//	// utility method
//	// TODO
//	insertMessage(gc.messages[:], message)
//}
//
//func (gc *Client) process() {
//	// process messages that have been sent down the various channels
//	for {
//		select {
//		case <-gc.shutdown:
//			return
//		case node := <-gc.aliveChan:
//			// merge the node in
//
//			gc.mergeNode(node) // TODO
//		case message := <-gc.messageChan:
//			// merge the message in
//
//			gc.mergeMessage(message) // TODO
//		}
//	}
//}

// replaced by sendNodes and sendMessages
//func (gc *Client) sendLoop() {
//	aliveTicker := time.NewTicker(gc.aliveTimeout)
//	messageTicker := time.NewTicker(gc.messageTimeout)
//
//loop:
//	for {
//		select {
//		case <-aliveTicker.C:
//			// send keepalive
//		case <-messageTicker.C:
//			// send message
//		case <-gc.shutdown:
//			break loop
//		}
//	}
//	aliveTicker.Stop()
//	messageTicker.Stop()
//	log.Println("send loop shut down")
//}

func (gc *Client) joinCluster(knownAddr net.IP) {
	// only ever called once, when you join the network
	node := newNodeDescriptor(knownAddr, time.Now(), 1, <-gc.counter.Count)
	//TODO: Save node? im not sure why knowaddr is used I thought we make a node of ourselves here and send it to knownAddr
	insertNode(gc.nodes[:], node)

	

}

// === exposed methods ===

func (gc *Client) Send(message string) error {
	// send a new message to the network
	gc.messageChan <- messageDescriptor{
		Descriptor: Descriptor{
			Count:     <-gc.counter.Count,
			ID:        gc.id,
			Timestamp: time.Now(),
		},
		Content: message,
	}

	return nil
}

func (gc *Client) logFile(logs File, payload nodeDescriptor) {
	//write information to txt file
	//display message on the CLI



}

func (gc *Clinet) createFile(){

}

func (gc *Client) Shutdown() {
	close(gc.shutdown)
}

func (gc *Client) Run(knownAddr net.IP) {
	gc.joinCluster(knownAddr)
	

	defer f.Close() //should the close be here?

	go gc.recvMessages()
	go gc.recvAlives()

	go gc.sendMessages()
	go gc.sendAlives()

	go gc.aliveLoop()
	go gc.messageLoop()

	//go gc.process()
}

// TODO:
//  - Make select random node/message thread-safe
