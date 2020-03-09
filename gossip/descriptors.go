package gossip

import (
	"hash/fnv"
	"net"
	"time"
)

// inherited by other xyzDescriptor types
type Descriptor struct {
	Timestamp time.Time
	ID        uint64
	Count     uint64
}

type nodeDescriptor struct {
	Descriptor
	Address *net.IP
}

func (d *Descriptor) time() time.Time {
	return d.Timestamp
}

func newNodeDescriptor(address *net.IP, timestamp time.Time, id uint64, count uint64) nodeDescriptor {
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
func insert(nodes []node, descriptor node, maxSize int) bool {
	oldest := -1 // index of oldest node
	for i, e := range nodes {
		if e.collisionHash() == descriptor.collisionHash() {
			// same, just update e
			// return true
		}
		if oldNode := nodes[oldest]; oldNode.time().After(e.time()) {
			// if current node is older than oldNode, set oldest node index to current index
			oldest = i
		}
	}
	if oldest == -1 {
		return false
	}
	if len(nodes) == maxSize {
		// replace
	} else {
		// add to array
	}
	return true
}

func merge(nodes1 []node, nodes2 []node, maxSize int) {
	// haha don't ask me about efficiency
	// NOTE: nodes2 max size will be nodes1 maxSize + 1, but it does not actually matter.
	for _, e := range nodes2 {
		insert(nodes1, e, maxSize)
	}
}
