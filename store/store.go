package store

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Record item in store
type Record struct {
	ID        string
	Data      json.RawMessage
	CreatedAt time.Time
}

var store map[string]Record

func init() {
	store = make(map[string]Record)
}

// GetItemsCount returns total count of elements
func GetItemsCount() (c int) {
	c = len(store)
	return
}

// GetRecord return record by ID
func GetRecord(id string) json.RawMessage {
	// @TODO process invalid request error
	rec, ok := store[id]
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
	store[id] = record
	return
}

// ReplaceRecord replace exists record by id
func ReplaceRecord(id string, data json.RawMessage) (ok bool) {
	rec, ok := store[id]
	if ok {
		rec.Data = data
		store[id] = rec
	}

	return
}

// DeleteRecord delete record by ID
func DeleteRecord(id string) bool {
	_, ok := store[id]
	if ok {
		delete(store, id)
	}
	return ok
}
