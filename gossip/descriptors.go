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
	Address net.IP
}

func (d *Descriptor) time() time.Time {
	return d.Timestamp
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
			nodes[i] = descriptor
			return true
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
			nodes[i] = descriptor
			return true
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