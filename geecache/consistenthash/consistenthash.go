package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)
//hash maps bytes to uint32
type Hash func(data []byte) uint32


type Map struct {
	hash Hash //hash方法
	replicas int //虚拟节点倍数
	keys []int //哈希环keys
	hashMap map[int]string //虚拟节点与真实节点映射表
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash: fn,
		replicas: replicas,
		hashMap: make(map[int]string),
	}
	//默认使用crc32的加密
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//添加真实节点/机器 的方法
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ { //设置虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

//针对key 计算所在节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx % len(m.keys)]]
}