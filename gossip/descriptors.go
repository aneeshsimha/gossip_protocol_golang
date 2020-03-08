package gossip

import (
	"hash/fnv"
	"net"
)

type nodeDescriptor struct {
	address   *net.IP
	timestamp uint64
	id        uint64
	addr      net.Addr
}

func (nd *nodeDescriptor) time() uint64 {
	return nd.timestamp
}

func (nd *nodeDescriptor) hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(nd.address.String()))
	return h.Sum64()
}

type messageDescriptor struct {
	content   []byte
	timestamp uint64
	id        uint64
	addr      net.Addr
}

func (md *messageDescriptor) time() uint64 {
	return md.timestamp
}

func (md *messageDescriptor) hash() uint64 {
	h := fnv.New64a()
	h.Write(md.content)
	return h.Sum64()
}

type node interface {
	time() uint64
	hash() uint64
}

//type nodeSet struct {
//	nodes   []node
//	maxSize int
//}

// TODO: wip
func insert(nodes []node, descriptor node, maxSize int) bool {
	oldest := -1
	for i, e := range nodes {
		if e.hash() == descriptor.hash() {
			// same, just update e
			// return true
		}
		if oldNode := nodes[oldest]; oldNode.time() > descriptor.time() {
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
