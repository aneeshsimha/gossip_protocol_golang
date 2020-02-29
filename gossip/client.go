package gossip

import (
	"log"
	"net/http"
)

type Client struct {
	nodes    []nodeDescriptor
	maxNodes uint64

	messages    []messageDescriptor
	maxMessages uint64

	shutdown chan bool

	port uint64
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
	// select loop to randomly select and send messages in the queue to known nodes
loop:
	// TODO: Not sure how I feel about using labels here; we could switch the break to an internal return...
	//  (See also: gc.sendAlives(), of course.)
	for {
		select {
		case <-gc.shutdown:
			// shutdown
			break loop
		default:
			// do something
		}
	}
	log.Printf("shutting down mssage send loop")
}

func (gc *Client) sendAlives() {
	// select loop to randomly select and send keepalives to nodes in the queue
loop:
	for {
		select {
		case <-gc.shutdown:
			// shutdown
			break loop
		default:
			// do something
		}
	}
	log.Printf("shutting down keep-alive loop")
}

func (gc *Client) Send(message string) error {
	// send a new message to the network

	return nil
}

func (gc *Client) Run() error {
	http.HandleFunc("/alive/", gc.aliveHandler)
	http.HandleFunc("/message/", gc.messageHandler)
	go gc.sendMessages()
	go gc.sendAlives()
	return http.ListenAndServe(":8080", nil)
}
