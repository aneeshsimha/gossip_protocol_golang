package gossip

import (
	"math/rand"
	"testing"
	"time"
)

var contents = [...]string{
	"hello world",
	"a sober possum is a social asset",
	"hell yeah brother cheers from iraq",
	"just vibing",
	"check out these lizards",
	"hey mr talk radio...",
	"traaaaaaains",
	"like a box of chocolates",
	"live. die. repeat.",
	"why does eve care so much about alice and bob",
}

func PrepNodes(now time.Time) []node {
	nodes := make([]node, len(contents))
	for i, _ := range nodes {
		nodes[i] = &messageDescriptor{
			Descriptor: Descriptor{
				Timestamp: now.Add(10 * time.Duration(i) * time.Second),
				ID:        rand.Uint64() % 5, // these don't matter so they can be truly random
				Count:     uint64(i),
			},
			Content: contents[i],
		}
	}
	return nodes
}

func printNodes(t *testing.T, nodes []node) {
	for _, e := range nodes {
		t.Log(e)
	}
}

func TestNodeMerge(t *testing.T) {
	rand.Seed(0) // same rand every time
	now := time.Now()
	nodes := PrepNodes(now)
	// "randomize"
	nodes[0], nodes[7] = nodes[7], nodes[0]
	nodes[3], nodes[5] = nodes[5], nodes[3]
	nodes[4], nodes[9] = nodes[9], nodes[4]
	nodes[1], nodes[2] = nodes[2], nodes[1]
	//rand.Shuffle(len(nodes), func(i, j int) {
	//	nodes[i], nodes[j] = nodes[j], nodes[i]
	//})

	printNodes(t, nodes)

	count := len(contents)
	moreNodes := []node{
		&messageDescriptor{
			Descriptor: Descriptor{
				Timestamp: now.Add(-10 * time.Second),
				ID:        rand.Uint64() % 5,
				Count:     uint64(count),
			},
			Content: "this shouldn't appear",
		},
		&messageDescriptor{
			Descriptor: Descriptor{
				Timestamp: now.Add(1000 * time.Second),
				ID:        rand.Uint64() % 5,
				Count:     uint64(count + 1),
			},
			Content: "this should replace hello world",
		},
	}
	count += 2

	//merge(nodes, moreNodes, len(nodes))

	if insert(nodes[:], moreNodes[0], len(nodes)) {
		t.Log("replaced something")
	} else {
		t.Log("NO REPLACE:", moreNodes[0], "was too old")
	}

	if insert(nodes[:], moreNodes[1], len(nodes)) {
		t.Log("replaced something")
	} else {
		t.Log("NO REPLACE:", moreNodes[1], "was too old")
	}

	printNodes(t, nodes)

	if nodes[7].collisionHash() != moreNodes[1].collisionHash() {
		t.Errorf("the wrong node was replaced")
	}
}
