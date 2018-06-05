package chash

import (
	"crypto/sha1"
	"sync"
	//	"hash"
	"math"
	"sort"
	"strconv"
)

const (
	//DefaultVirualSpots default virual spots
	DefaultVirualSpots = 400
)

type node struct {
	nodeKey   string
	spotValue uint32
}

type nodesArray []node

func (p nodesArray) Len() int           { return len(p) }
func (p nodesArray) Less(i, j int) bool { return p[i].spotValue < p[j].spotValue }
func (p nodesArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p nodesArray) Sort()              { sort.Sort(p) }

//HashRing2 store nodes and weigths
type HashRing2 struct {
	virualSpots int
	nodes       nodesArray
	weights     map[string]int
	mu          sync.RWMutex
}

//NewHashRing2 create a hash ring with virual spots
func NewHashRing2(spots int) *HashRing2 {
	if spots == 0 {
		spots = DefaultVirualSpots
	}

	h := &HashRing2{
		virualSpots: spots,
		weights:     make(map[string]int),
	}
	return h
}

//AddNodes add nodes to hash ring
func (h *HashRing2) AddNodes(nodeWeight map[string]int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for nodeKey, w := range nodeWeight {
		h.weights[nodeKey] = w
	}
	h.generate()
}

//AddNode add node to hash ring
func (h *HashRing2) AddNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

//RemoveNode remove node
func (h *HashRing2) RemoveNode(nodeKey string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.weights, nodeKey)
	h.generate()
}

//UpdateNode update node with weight
func (h *HashRing2) UpdateNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

func (h *HashRing2) generate() {
	var totalW int
	for _, w := range h.weights {
		totalW += w
	}

	totalVirtualSpots := h.virualSpots * len(h.weights)
	h.nodes = nodesArray{}

	for nodeKey, w := range h.weights {
		spots := int(math.Floor(float64(w) / float64(totalW) * float64(totalVirtualSpots)))
		for i := 1; i <= spots; i++ {
			hash := sha1.New()
			hash.Write([]byte(nodeKey + ":" + strconv.Itoa(i)))
			hashBytes := hash.Sum(nil)
			n := node{
				nodeKey:   nodeKey,
				spotValue: genValue(hashBytes[6:10]),
			}
			h.nodes = append(h.nodes, n)
			hash.Reset()
		}
	}
	h.nodes.Sort()
}

func genValue(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
	return v
}

//GetNode get node with key
func (h *HashRing2) GetNode(s string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.nodes) == 0 {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)
	v := genValue(hashBytes[6:10])
	i := sort.Search(len(h.nodes), func(i int) bool { return h.nodes[i].spotValue >= v })

	if i == len(h.nodes) {
		i = 0
	}
	return h.nodes[i].nodeKey
}

// func main() {
// 	const (
// 		node1 = "192.168.1.1"
// 		node2 = "192.168.1.2"
// 		node3 = "192.168.1.3"
// 		node4 = "192.168.1.4"
// 		node5 = "192.168.1.5"
// 	)
// 	nodeWeight := make(map[string]int)
// 	nodeWeight[node1] = 1
// 	nodeWeight[node2] = 1
// 	nodeWeight[node3] = 1
// 	nodeWeight[node4] = 1
// 	nodeWeight[node5] = 1
// 	vitualSpots := 2

// 	hash := NewHashRing2(vitualSpots)
// 	hash.AddNodes(nodeWeight)

// 	fmt.Println(hash.nodes)
// 	ipMap := make(map[string]int, 0)
// 	for i := 0; i < 1000; i++ {
// 		si := fmt.Sprintf("key%d", i)
// 		k := hash.GetNode(si)
// 		if _, ok := ipMap[k]; ok {
// 			ipMap[k]++
// 		} else {
// 			ipMap[k] = 1
// 		}
// 	}

// 	for k, v := range ipMap {
// 		fmt.Println("Node IP:", k, " count:", v)
// 	}
// }
