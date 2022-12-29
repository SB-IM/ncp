package cache

// Fork From: https://leetcode-cn.com/problems/lru-cache/solution/lruhuan-cun-ji-zhi-by-leetcode-solution/

type LRUCache struct {
	capacity   int
	cache      map[string]*DLinkedNode
	head, tail *DLinkedNode
}

type DLinkedNode struct {
	key, value string
	prev, next *DLinkedNode
}

func initDLinkedNode(key, value string) *DLinkedNode {
	return &DLinkedNode{
		key:   key,
		value: value,
	}
}

func NewLRUCache(capacity int) LRUCache {
	l := LRUCache{
		cache:    map[string]*DLinkedNode{},
		head:     initDLinkedNode("", ""),
		tail:     initDLinkedNode("", ""),
		capacity: capacity,
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

func (this *LRUCache) Get(key string) string {
	if _, ok := this.cache[key]; !ok {
		return ""
	}
	node := this.cache[key]
	this.moveToHead(node)
	return node.value
}

func (this *LRUCache) Put(key string, value string) {
	if _, ok := this.cache[key]; !ok {
		node := initDLinkedNode(key, value)
		this.cache[key] = node
		this.addToHead(node)
		if len(this.cache) > this.capacity {
			removed := this.removeTail()
			delete(this.cache, removed.key)
		}
	} else {
		node := this.cache[key]
		node.value = value
		this.moveToHead(node)
	}
}

func (this *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = this.head
	node.next = this.head.next
	this.head.next.prev = node
	this.head.next = node
}

func (this *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (this *LRUCache) moveToHead(node *DLinkedNode) {
	this.removeNode(node)
	this.addToHead(node)
}

func (this *LRUCache) removeTail() *DLinkedNode {
	node := this.tail.prev
	this.removeNode(node)
	return node
}
