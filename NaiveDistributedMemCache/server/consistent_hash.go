package server

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashRing struct {
	keys []int
	virtualReplicas int
	virtualNodeMap map[int]string
	hash func(data []byte) uint32
}

func NewHashRing(virtualReplicas int, hash func(data []byte) uint32) *HashRing {
	hashRing := &HashRing{
		virtualReplicas: virtualReplicas,
		virtualNodeMap:  make(map[int]string),
		hash:            hash,
	}
	if hash == nil {
		hashRing.hash = crc32.ChecksumIEEE
	}
	return hashRing
}

func (ring *HashRing) AddNodes(keys ...string) {
	for _, key := range keys {
		for i := 0; i < ring.virtualReplicas; i++ {
			hash := int(ring.hash([]byte(strconv.Itoa(i) + key)))
			ring.keys = append(ring.keys, hash)
			ring.virtualNodeMap[hash] = key
		}
	}
	sort.Ints(ring.keys)
}

func (ring *HashRing) GetNode(key string) string {
	if len(ring.keys) == 0 {
		return ""
	}
	hash := int(ring.hash([]byte(key)))
	idx := sort.Search(len(ring.keys), func(i int) bool {
		return ring.keys[i] >= hash
	})
	return ring.virtualNodeMap[ring.keys[idx%(len(ring.keys))]]
}
