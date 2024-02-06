package selector

import (
	"hash/crc32"
	"rpc-oneway/pkg/resolver"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type HashRing struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]*resolver.ServiceInstance
}

func NewHashRing(replicas int, fn Hash) *HashRing {
	m := &HashRing{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]*resolver.ServiceInstance),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty returns true if there are no items available.
func (m *HashRing) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add adds some keys to the hash.
func (m *HashRing) Add(ins ...resolver.ServiceInstance) {
	for _, key := range ins {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key.Endpoint))) //计算key的hash
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = &key
		}
	}

	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *HashRing) Get(key string) *resolver.ServiceInstance {
	if m.IsEmpty() {
		return nil
	}

	hash := int(m.hash([]byte(key)))

	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash }) //找到最近的合适的node

	// Means we have cycled back to the first replica.
	if idx == len(m.keys) {
		idx = 0
	}

	//fmt.Printf("LINK: key(%s), keyHash(%d), mapLen(%d), keysLen(%d), idx(%d), node(%s), keys(%v), nodes(%v) \n",
	//	key, hash, len(m.hashMap), len(m.keys), idx, m.hashMap[m.keys[idx]].node, m.keys, m.hashMap)

	return m.hashMap[m.keys[idx]]
}
