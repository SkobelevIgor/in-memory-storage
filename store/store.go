package store

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Record is basic struct for object in store
type Record struct {
	ID        string
	Data      interface{}
	CreatedAt time.Time
}

var store map[string]Record

func init() {
	store = make(map[string]Record)
}

// GetRecord return record by ID
func GetRecord(id string) (record interface{}, err error) {
	// @TODO process invalid request error
	record, ok := store[id]
	a := Record{ID: "'one'"}
	if ok {
		return a, nil
	}
	return a, errors.New("Record not found")
}

// SaveRecord save new record to store
func SaveRecord(data interface{}) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	id := uid.String()
	record := Record{ID: id, Data: data, CreatedAt: time.Now()}
	store[id] = record
	return id, err
}

// UpdateRecord update exists record by id
// @TODO recheck data
// @TODO recheck REST return on update
// func ReplaceRecord(id int64, data interface{}) interface{}, error  {
// 	return
// }

// // DeleteRecord delete record by ID
// func DeleteRecord(id int64) int64 {
// 	return id
// }
