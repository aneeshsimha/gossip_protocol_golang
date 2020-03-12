package gossip

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aneeshsimha/gossip_protocol_golang/counter"
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
	rand.Seed(time.Now().UnixNano())
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

func (gc *Client) sendAlives() {
	// select a random node Descriptor and send keepalive

	aliveTicker := time.NewTicker(gc.aliveTimeout)
	defer aliveTicker.Stop()
	defer log.Println("sendAlives loop shut down")

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

func (gc *Client) sendMessages() {
	// 1. select a random message and node Descriptor
	// 2. send message to node, and request a message from node
	// 3. merge reply into own message slice

	// select a random node Descriptor and send random message

	messageTicker := time.NewTicker(gc.messageTimeout)
	defer messageTicker.Stop()
	defer log.Println("sendMessages loop shut down")

	// loop forever
	for {
		select {
		case <-messageTicker.C: // do every interval
			//  choose a random known node descriptor
			//  choose a random stored message
			//  turn the messageDescriptor into a stringPacket
			//  send the stringPacket
			gc.sendMessage()
		case <-gc.shutdown:
			return
		}
	}
}

func (gc *Client) sendAlive() {
	// choose a random known node descriptor
	// add own ip + current time to the copy of the node Descriptor list before sending

	randomNode := randomNeighbor(gc.nodes)
	if randomNode.Address == nil {
		return
	}

	conn, err := net.Dial("tcp", randomNode.Address.String()+":"+gc.alivePort)
	if err != nil {
		return
	}
	defer conn.Close() // this is why this is its own function

	// make/get self node
	me := newNodeDescriptor(nil, time.Now(), gc.id, <-gc.counter.Count) // filled in on the other end
	kap := prepareKeepAlivePayload(gc.nodes, me)

	enc := gob.NewEncoder(conn)
	_ = enc.Encode(kap)
}

func (gc *Client) sendMessage() {
	// choose random node descriptor
	// choose random message
	// send message

	// choose random node
	randomNode := randomNeighbor(gc.nodes)
	if randomNode.Address == nil {
		return
	}

	// choose random message
	randMessage := randomMessage(gc.messages)
	if randMessage.Content == "" {
		return
	}

	// TCP connection to the random node
	conn, err := net.Dial("tcp", randomNode.Address.String()+":"+gc.messagePort)
	if err != nil {
		return
	}
	defer conn.Close() // this is why this is its own function
	//create payload to send message
	payload := newStringPayload(randMessage)

	enc := gob.NewEncoder(conn)
	_ = enc.Encode(payload)
}

// Channel used to recieve keepAlive packets, concurrency provided by spinning off a new thread each time a connection is established 
func (gc *Client) recvAlives() {
	listener, err := net.Listen("tcp", ":"+gc.alivePort)
	defer log.Println("shut down recvAlives")
	if err != nil {
		log.Fatal("could not listen for keep alives:", err)
	}
	for {
		select {
		case <-gc.shutdown:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go gc.handleAlive(conn)
		}
	}
}

// same as recvAlive except it handles messages
func (gc *Client) recvMessages() {
	listener, err := net.Listen("tcp", ":"+gc.messagePort)
	if err != nil {
		log.Fatal("could not listen for messages:", err)
	}
	for {
		select {
		case <-gc.shutdown:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go gc.handleMessage(conn)
		}
	}
}

// decodes node packet and inserts each desc into the client's partial view, only keeping the newest ones
func (gc *Client) handleAlive(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	kap := KeepAlivePayload{}
	dec.Decode(&kap)

	for _, desc := range kap.KnownNodes {
		if desc.ID == gc.id {
			continue // don't bother adding own originating messages
		}
		if desc.Address == nil {
			// if it's a nil address, then it's the address of the sender
			strAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
			desc.Address = net.ParseIP(strAddr)
		}
		gc.logFile(desc)
		gc.aliveChan <- desc
	}
}

//same functionality as handleAlive, except there is only one message per request
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

// channel that inserts a given node desc into the client's view, keeping newest descs first
func (gc *Client) aliveLoop() {
	defer log.Println("aliveLoop shut down")
	for {
		select {
		case <-gc.shutdown:
			return
		case desc := <-gc.aliveChan:
			insertNode(gc.nodes[:], desc)
		}
	}
}

// channel inserts message into message array
func (gc *Client) messageLoop() {
	defer log.Println("messageLoop shut down")
	for {
		select {
		case <-gc.shutdown:
			return
		case desc := <-gc.messageChan:
			insertMessage(gc.messages[:], desc)
		}
	}
}

// Used at startup to join an existing node cluster or to start a node cluster. Needs the ip address of at least 1 running node
func (gc *Client) joinCluster(knownAddr net.IP) {
	// only ever called once, when you join the network
	node := newNodeDescriptor(knownAddr, time.Now(), 1, <-gc.counter.Count)
	insertNode(gc.nodes[:], node)

}

// === exposed methods ===

func (gc *Client) Send(message string) error {
	// send a new message to the network
	msg := messageDescriptor{
		Descriptor: Descriptor{
			Count:     <-gc.counter.Count,
			ID:        gc.id,
			Timestamp: time.Now(),
		},
		Content: message,
	}
	gc.messageChan <- msg

	log.Println("created message:", msg)

	return nil
}

func (gc *Client) logFile(payload nodeDescriptor) {
	//write information to txt file
	//display message on the CLI

	f, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bytes := []byte(fmt.Sprintf(
		"ID of Origin: %d\t Node Address: %s\t Time: %v\tappended some data\n",
		payload.ID,
		payload.Address.String(),
		payload.Timestamp,
	))
	if _, err := f.Write(bytes); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (gc *Client) Shutdown() {
	close(gc.shutdown)
}

func (gc *Client) Run(knownAddr string) {
	// a nil addr means it's the first node, others will join
	if knownAddr != "" {
		gc.joinCluster(net.ParseIP(knownAddr))
	}
	go gc.recvMessages()
	go gc.recvAlives()

	go gc.sendMessages()
	go gc.sendAlives()

	go gc.aliveLoop()
	go gc.messageLoop()
}

func (gc *Client) Nodes() []nodeDescriptor {
	return gc.nodes
}

func (gc *Client) Messages() []messageDescriptor {
	return gc.messages
}
