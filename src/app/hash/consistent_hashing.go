package chash

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

const DEFAULT_REPLICAS = 500

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Node struct {
	Id         int
	Ip         string
	Port       int
	HostName   string
	Weight     int
	Collection string
}

func NewNode(id int, ip string, port int, name string, weight int, coll string) *Node {
	return &Node{
		Id:         id,
		Ip:         ip,
		Port:       port,
		HostName:   name,
		Weight:     weight,
		Collection: coll,
	}
}

type Consistent struct {
	Nodes     map[uint32]Node
	numReps   int
	Resources map[int]bool
	ring      HashRing
	sync.RWMutex
}

func NewConsistent(replicas int) *Consistent {
	if replicas == 0 {
		replicas = DEFAULT_REPLICAS
	}
	return &Consistent{
		Nodes:     make(map[uint32]Node),
		numReps:   replicas,
		Resources: make(map[int]bool),
		ring:      HashRing{},
	}
}

func (c *Consistent) Add(node *Node) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Id]; ok {
		return false
	}

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		c.Nodes[c.hashStr(str)] = *(node)
	}
	c.Resources[node.Id] = true
	c.sortHashRing()
	return true
}

func (c *Consistent) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

func (c *Consistent) joinStr(i int, node *Node) string {
	return node.Ip + "*" + strconv.Itoa(node.Weight) +
		"-" + strconv.Itoa(i) +
		"-" + strconv.Itoa(node.Id)
}

func (c *Consistent) hashStr(key string) uint32 {
	// return crc32.ChecksumIEEE([]byte(key))

	// md5Chan := make(chan []byte, 1)
	// md5Sum := md5.Sum([]byte(key))
	// md5Chan <- md5Sum[:]
	// return crc32.ChecksumIEEE(<-md5Chan)

	hash := sha1.New()
	hash.Write([]byte(key))
	bs := hash.Sum(nil)
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
	return v
}

func (c *Consistent) Get(key string) Node {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashStr(key)
	i := c.search(hash)

	return c.Nodes[c.ring[i]]
}

func (c *Consistent) search(hash uint32) int {

	i := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if i < len(c.ring) {
		return i
	}
	return 0
}

func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Id]; !ok {
		return
	}

	delete(c.Resources, node.Id)

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	c.sortHashRing()
}

func main() {

	cHashRing := NewConsistent(1000)

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		cHashRing.Add(NewNode(i, "172.18.1."+si, 8080, "host_"+si, 1, ""))
	}

	ipMap := make(map[string]int, 0)
	for i := 0; i < 100000; i++ {
		si := fmt.Sprintf("key%d", i)
		k := cHashRing.Get(si)
		if _, ok := ipMap[k.Ip]; ok {
			ipMap[k.Ip]++
		} else {
			ipMap[k.Ip] = 1
		}
	}

	for k, v := range ipMap {
		fmt.Println("Node IP:", k, " count:", v)
	}

}
