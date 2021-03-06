package gossip

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"
	"time"
)

// useful utility functions for descriptors
// inherited by other xyzDescriptor types
type Descriptor struct {
	Timestamp time.Time
	ID        uint64 // id of originator gossip client
	Count     uint64 // count value, shouldn't ever be repeated; (ID, Count) tuple for a unique identifying pair
}

func (d Descriptor) String() string {
	return fmt.Sprintf("{Timestamp: %v, ID: %v, Count: %v}", d.Timestamp, d.ID, d.Count)
}

func (d *Descriptor) time() time.Time {
	return d.Timestamp
}

type nodeDescriptor struct {
	Descriptor
	Address net.IP
}

func (nd nodeDescriptor) String() string {
	return fmt.Sprintf("{nodeDescriptor:: %s, Address: %s}", nd.Descriptor, nd.Address)
}

func newNodeDescriptor(address net.IP, timestamp time.Time, id uint64, count uint64) nodeDescriptor {
	return nodeDescriptor{
		Descriptor: Descriptor{
			Timestamp: timestamp,
			ID:        id,
			Count:     count,
		},
		Address: address,
	}
}

func (nd *nodeDescriptor) collisionHash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(nd.Address.String()))
	return h.Sum64()
}

type messageDescriptor struct {
	Descriptor
	Content string
}

func (md messageDescriptor) String() string {
	return fmt.Sprintf("{messageDescriptor:: %s, Content: %s}", md.Descriptor, md.Content)
}

func newMessageDescriptor(content string, timestamp time.Time, id uint64, count uint64) messageDescriptor {
	return messageDescriptor{
		Descriptor: Descriptor{
			Timestamp: timestamp,
			ID:        id,
			Count:     count,
		},
		Content: content,
	}
}

func (md *messageDescriptor) collisionHash() uint64 {
	h := fnv.New64a()
	hashStr := fmt.Sprintf("(%d, %d)", md.ID, md.Count)
	h.Write([]byte(hashStr))
	return h.Sum64()
}

type node interface {
	time() time.Time
	collisionHash() uint64
}

func insertNode(nodes []nodeDescriptor, descriptor nodeDescriptor) bool {
	oldest := 0 // index of oldest node
	for i, e := range nodes {
		if e.ID == 0 {
			// e is uninitialized
			// array is not full, so just "append"
			oldest = i
			break
		}
		if e.collisionHash() == descriptor.collisionHash() {
			// node is already in the table, just update e
			if e.time().Before(descriptor.time()) {
				nodes[i] = descriptor
				return true
			}
			return false

		}

		// else, check if the current node is older than the current old node
		if oldNode := nodes[oldest]; oldNode.time().After(e.time()) {
			// if current node is older than oldNode, set oldest node index to current index
			oldest = i
		}
	}
	if nodes[oldest].time().After(descriptor.time()) {
		// descriptor is older than the oldest node
		return false
	}
	nodes[oldest] = descriptor
	return true
}

func insertMessage(messages []messageDescriptor, descriptor messageDescriptor) bool {
	oldest := 0 // index of oldest node
	for i, e := range messages {
		if e.ID == 0 {
			// e is uninitialized
			// array is not full, so just "append"
			oldest = i
			break
		}
		if e.collisionHash() == descriptor.collisionHash() {
			// node is already in the table, just update e
			if e.time().Before(descriptor.time()) {
				messages[i] = descriptor
				return true
			}
			return false
		}

		// else, check if the current node is older than the current old node
		if oldNode := messages[oldest]; oldNode.time().After(e.time()) {
			// if current node is older than oldNode, set oldest node index to current index
			oldest = i
		}
	}
	if messages[oldest].time().After(descriptor.time()) {
		// descriptor is older than the oldest node
		return false
	}
	fmt.Println("inserting message:", descriptor)
	messages[oldest] = descriptor
	return true
}

// utility functions to select random node
// tries to get a non-nil node 5x before failing
func randomNeighbor(nodes1 []nodeDescriptor) nodeDescriptor {
	randomNode := nodes1[rand.Intn(len(nodes1))]
	for i := 0; i < 5; i += 1 {
		// try to get an existing random node 5 times
		if randomNode.Address != nil {
			break
		} else {
			randomNode = nodes1[rand.Intn(len(nodes1))]
		}
	}
	return randomNode
}

func randomMessage(messages []messageDescriptor) messageDescriptor {
	randomMsg := messages[rand.Intn(len(messages))]
	for i := 0; i < 5; i += 1 {
		// try to get an existing random node 5 times
		if randomMsg.Content != "" {
			break
		} else {
			randomMsg = messages[rand.Intn(len(messages))]
		}
	}
	return randomMsg
}
