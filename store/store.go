package store

import (
	"encoding/json"
	"errors"
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
func SaveRecord(data json.RawMessage) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	id := uid.String()
	record := Record{ID: id, Data: data, CreatedAt: time.Now()}
	store[id] = record
	return id, err
}

// ReplaceRecord replace exists record by id
func ReplaceRecord(id string, data json.RawMessage) (err error) {
	rec, ok := store[id]
	if ok {
		rec.Data = data
		store[id] = rec
	} else {
		err = errors.New("Unable to find record")
	}

	return
}

// DeleteRecord delete record by ID
func DeleteRecord(id string) (err error) {
	_, ok := store[id]
	if ok {
		delete(store, id)
	}
	return
}
