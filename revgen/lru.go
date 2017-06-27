// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"container/list"
	"encoding/json"
	"log"
	"sync"
)

const defaultCapacity = 1000

type UserRecord struct {

	// unique per POST request (job id)
	Uuid string

	RevBytes []byte
}

// A LRU cache that is safe for concurrent access. All items are also added to
// persistent storage.
type LruCache struct {
	sync.RWMutex

	// Items are evicted when cache is full, can be fetched from persistent storage
	Capacity int

	// The Value of each list element will be a *UserRecord
	linkedList *list.List

	// The Value of each list element will be a *UserRecord
	catalog map[string]*list.Element

	// Persistent storage for results by uuid
	dao *UserAccountFile
}

func NewCache(capacity int, filename string) *LruCache {

	useCapacity := defaultCapacity
	if capacity > 0 {
		useCapacity = capacity
	}
	return &LruCache{
		Capacity:   useCapacity,
		linkedList: list.New(),
		catalog:    make(map[string]*list.Element, useCapacity), // size hint is capacity
		dao:        New(filename),
	}
}

// Assumes record is not already in cache, which seems reasonable given that
// the key is uuid. Consider adding lookup for safety.
func (cache *LruCache) Add(uuid string, revbytes []byte) {

	if cache == nil {
		return
	}
	if len(uuid) == 0 {
		return
	}
	log.Printf("Adding reversed bytes to cache for job %q\n", uuid)
	// Takes a write lock while mutating cache
	cache.Lock()
	defer cache.Unlock()

	r := UserRecord{uuid, revbytes}
	element := cache.linkedList.PushFront(&r)
	cache.catalog[uuid] = element
	// Persist to disk so that eviction does not cause data loss
	// TODO: check result for errors
	cache.dao.Append(r)
	if cache.linkedList.Len() > cache.Capacity {
		go cache.Evict()
	}
}

// Purge the oldest entry
func (cache *LruCache) Evict() {
	if cache == nil {
		return
	}
	cache.Lock()
	defer cache.Unlock()
	element := cache.linkedList.Back()
	if element != nil {
		uuid := element.Value.(*UserRecord).Uuid
		delete(cache.catalog, uuid)
		cache.linkedList.Remove(element)
	}
}

func (cache *LruCache) Get(uuid string) (value *UserRecord, ok bool) {

	if cache.catalog == nil {
		return
	}
	log.Printf("getting record with uuid: %s\n", uuid)
	cache.RLock()
	element, ok := cache.catalog[uuid]
	cache.RUnlock()
	if ok {
		cache.Lock()
		defer cache.Unlock()
		cache.linkedList.MoveToFront(element)
		if m, ok := element.Value.(*UserRecord); ok {
			return m, true
		}
	} else {
		// TODO: Seems expected that this should Add to cache
		if b, ok := cache.dao.Get(uuid); ok {
			var m UserRecord
			if err := json.Unmarshal(b, &m); err == nil {
				return &m, ok
			}
		}
	}
	return
}
