package gossip

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"
	"time"
)

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

//func (nd *nodeDescriptor) time() uint64 {
//	return nd.Timestamp
//}

func (nd *nodeDescriptor) collisionHash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(nd.Address.String()))
	return h.Sum64()
}

type messageDescriptor struct {
	Descriptor
	Content string
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

//func (md *messageDescriptor) time() uint64 {
//	return md.Timestamp
//}

func (md *messageDescriptor) collisionHash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(md.Content))
	return h.Sum64()
}

type node interface {
	time() time.Time
	collisionHash() uint64
}

//type nodeSet struct {
//	nodes   []node
//	maxSize int
//}

// TODO: wip
//  Okay, so this doesn't work, because golang copies interface slices weirdly (linear time rather than constant)
//func insert(nodes []node, descriptor node, maxSize int) bool {
//	oldest := 0 // index of oldest node
//	for i, e := range nodes {
//		if e == nil {
//			// array is not full, so just "append"
//			oldest = i
//			break
//		}
//		if e.collisionHash() == descriptor.collisionHash() {
//			// node is already in the table, just update e
//			nodes[i] = descriptor
//			return true
//		}
//
//		// else, check if the current node is older than the current old node
//		if oldNode := nodes[oldest]; oldNode.time().After(e.time()) {
//			// if current node is older than oldNode, set oldest node index to current index
//			oldest = i
//		}
//	}
//	if nodes[oldest].time().After(descriptor.time()) {
//		// descriptor is older than the oldest node
//		return false
//	}
//	nodes[oldest] = descriptor
//	return true
//}

//func merge(nodes1 []node, nodes2 []node, maxSize int) {
//	// haha don't ask me about efficiency
//	// NOTE: nodes2 max size will be nodes1 maxSize + 1, but it does not actually matter.
//	for _, e := range nodes2 {
//		insert(nodes1, e, maxSize)
//	}
//}

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

func insertMessage(nodes []messageDescriptor, descriptor messageDescriptor) bool {
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

// utility functions to select random node
func randomNeighbor(nodes1 []nodeDescriptor) nodeDescriptor {
	iter := rand.Intn(len(nodes1))
	return nodes1[iter]
}

func randomMessage(nodes1 []messageDescriptor) messageDescriptor {
	iter := rand.Intn(len(nodes1))
	return nodes1[iter]
}
