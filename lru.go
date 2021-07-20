package rocinante

import "sync"

type LRUCache struct {
	keys  []string
	cache *sync.Map
	rw    *sync.RWMutex
	cap   int
}

func newCache(cap int) *LRUCache {
	return &LRUCache{
		keys:  make([]string, 0),
		cache: &sync.Map{},
		rw:    &sync.RWMutex{},
		cap:   cap,
	}
}

func defaultCache() *LRUCache {
	return &LRUCache{
		keys:  make([]string, 0),
		cache: &sync.Map{},
		rw:    &sync.RWMutex{},
		cap:   200,
	}
}

func (l *LRUCache) Get(key string) (interface{}, bool) {
	if !l.isHead(key) {
		l.rw.Lock()

		l.moveKeyToHead(key)

		l.rw.Unlock()
	}
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.cache.Load(key)
}

func (l *LRUCache) Set(key string, value interface{}) {
	l.rw.Lock()
	defer l.rw.Unlock()

	if _, exists := l.cache.Load(key); exists {
		l.moveKeyToHead(key)
		l.cache.Store(key, value)
		return
	}

	if len(l.keys) >= l.cap {
		lastKay := l.removeLastKey()
		l.cache.Delete(lastKay)
	}
	l.insertKeyToHead(key)
	l.cache.Store(key, value)
}

func (l *LRUCache) moveKeyToHead(key string) {
	for i, s := range l.keys {
		if s == key {
			tmp1 := l.keys[0:i]
			tmp2 := l.keys[i+1:]
			l.keys = append(tmp1, tmp2...)
			l.keys = append(l.keys, key)
		}
	}
}

func (l *LRUCache) insertKeyToHead(key string) {
	l.keys = append(l.keys, key)
}

func (l *LRUCache) removeLastKey() string {
	lastKey := l.keys[0]
	l.keys = l.keys[1:]
	return lastKey
}

func (l *LRUCache) isHead(key string) bool {
	if len(l.keys) > 0 {
		return l.keys[len(l.keys)-1] == key
	}
	return false
}
