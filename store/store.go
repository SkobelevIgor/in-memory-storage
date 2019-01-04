package store

import (
	"encoding/json"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Record item in store
type Record struct {
	ID        string
	Data      json.RawMessage
	CreatedAt time.Time
}

type mainStore struct {
	mx sync.RWMutex
	s  map[string]Record
}

var store *mainStore

func init() {
	store = &mainStore{s: make(map[string]Record)}
}

// GetItemsCount returns total count of elements
func GetItemsCount() (c int) {
	store.mx.RLock()
	defer store.mx.RUnlock()
	c = len(store.s)
	return
}

// GetRecord return record by ID
func GetRecord(id string) json.RawMessage {
	store.mx.RLock()
	defer store.mx.RUnlock()
	rec, ok := store.s[id]
	if ok {
		return rec.Data
	}
	return nil
}

// SaveRecord save new record to store
func SaveRecord(data json.RawMessage) (id string) {
	uid, err := uuid.NewV4()
	if err != nil {
		return
	}
	id = uid.String()
	record := Record{ID: id, Data: data, CreatedAt: time.Now()}
	store.mx.Lock()
	defer store.mx.Unlock()
	store.s[id] = record
	return
}

// ReplaceRecord replace exists record by id
func ReplaceRecord(id string, data json.RawMessage) (ok bool) {
	store.mx.Lock()
	defer store.mx.Unlock()
	rec, ok := store.s[id]
	if ok {
		rec.Data = data
		store.s[id] = rec
	}
	return
}

// DeleteRecord delete record by ID
func DeleteRecord(id string) bool {
	store.mx.Lock()
	defer store.mx.Unlock()
	_, ok := store.s[id]
	if ok {
		delete(store.s, id)
	}
	return ok
}
