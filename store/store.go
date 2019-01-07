package store

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Config store config
type Config struct {
	BackupTimeSecs     int
	BackupOps          int
	SkipGracefulBackup bool
	SkipBackupLoad     bool
}

// Record item in store
type Record struct {
	ID        string
	Data      json.RawMessage
	CreatedAt time.Time
}

type mainStore struct {
	mx              sync.RWMutex
	backupInProcess bool
	opsCount        int
	records         map[string]Record
}

var (
	store    *mainStore
	stopChan chan os.Signal
)

// InitStore initialize store, load from backup etc
func InitStore(cfg Config) (err error) {
	store = &mainStore{records: make(map[string]Record)}

	runBackupByTicker(cfg.BackupTimeSecs)
	opsCountCfg = cfg.BackupOps

	if cfg.SkipGracefulBackup == false {
		runStopListener()
	}
	if cfg.SkipBackupLoad == false {
		err = loadStoreFromBackup()
	}
	return err
}

func runStopListener() {
	stopChan = make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM)
	signal.Notify(stopChan, syscall.SIGINT)
	go func() {
		signal := <-stopChan
		fmt.Printf("Signal: %+v\n", signal)
		err := backupStore()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()
}

// GetItemsCount returns total count of elements
func GetItemsCount() int {
	store.mx.RLock()
	defer store.mx.RUnlock()
	return len(store.records)

}

// GetRecord return record by ID
func GetRecord(id string) json.RawMessage {
	store.mx.RLock()
	defer store.mx.RUnlock()
	rec, ok := store.records[id]
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
	store.records[id] = record
	store.opsCount++
	store.mx.Unlock()
	runBackupByOpsCounter()
	return
}

// ReplaceRecord replace exists record by id
func ReplaceRecord(id string, data json.RawMessage) (ok bool) {
	store.mx.Lock()
	rec, ok := store.records[id]
	if ok {
		rec.Data = data
		store.records[id] = rec
	}
	store.opsCount++
	store.mx.Unlock()
	runBackupByOpsCounter()
	return
}

// DeleteRecord delete record by ID
func DeleteRecord(id string) bool {
	store.mx.Lock()
	_, ok := store.records[id]
	if ok {
		delete(store.records, id)
	}
	store.opsCount++
	store.mx.Unlock()
	runBackupByOpsCounter()
	return ok
}
