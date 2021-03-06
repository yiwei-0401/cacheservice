package lru

import (
	"container/list"
)

type Cache struct {
	maxBtyes int64 //最大容量
	nbytes int64 //当前使用
	ll *list.List //双向链表
	cache map[string]*list.Element //map
	Onevicted func(key string, value Value) //回调
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBtyes: maxBytes,
		ll : list.New(),
		cache : make(map[string]*list.Element),
		Onevicted : onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	//命中后把节点移到队尾， 返回值
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		//类型断言
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	//删除队首节点，删掉map里的key，重新计算内存，调用回调函数
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.Onevicted != nil {
			c.Onevicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value)  {
	//有用的移到前面
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBtyes != 0 && c.maxBtyes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
