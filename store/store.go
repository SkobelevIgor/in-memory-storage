package store

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	backupFile     = ".store.gob"
	backupFileTemp = ".store.go.TMP"
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
	mx      sync.RWMutex
	records map[string]Record
}

var (
	store    *mainStore
	stopChan chan os.Signal
)

// InitStore initialize store, load from backup etc
func InitStore(cfg Config) (err error) {
	store = &mainStore{records: make(map[string]Record)}
	if cfg.BackupTimeSecs > 0 {
		err = runBackupByTimer(cfg.BackupTimeSecs)
	}
	if err == nil && cfg.BackupOps > 0 {
		err = runBackupByOpsCounter(cfg.BackupOps)
	}
	if err == nil && cfg.SkipGracefulBackup == false {
		runStopListener()
	}
	if err == nil && cfg.SkipBackupLoad == false {
		err = loadStoreFromBackup()
	}
	return err
}

func runBackupByTimer(secsInterval int) error {
	return nil
}

func runBackupByOpsCounter(opsCount int) error {
	return nil
}

func loadStoreFromBackup() (err error) {
	if _, err := os.Stat(backupFile); !os.IsNotExist(err) {
		f, err := os.Open(backupFile)
		if err == nil {
			decoder := gob.NewDecoder(f)
			recs := map[string]Record{}
			err = decoder.Decode(&recs)
			store.records = recs
		}
	}
	return err
}

func backupStore() (err error) {
	fmt.Println("Saving store to disk ...")
	f, err := os.Create(backupFileTemp)
	if err == nil {
		defer f.Close()
		enc := gob.NewEncoder(f)
		enc.Encode(store.records)
		err = os.Rename(backupFileTemp, backupFile)
		fmt.Println("Store saved to disk")
	}
	return err
}

func runStopListener() {
	stopChan = make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM)
	signal.Notify(stopChan, syscall.SIGINT)
	go func() {
		signal := <-stopChan
		fmt.Printf("caught sig: %+v\n", signal)
		err := backupStore()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()
}

// GetItemsCount returns total count of elements
func GetItemsCount() (c int) {
	store.mx.RLock()
	defer store.mx.RUnlock()
	c = len(store.records)
	return
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
	defer store.mx.Unlock()
	store.records[id] = record
	return
}

// ReplaceRecord replace exists record by id
func ReplaceRecord(id string, data json.RawMessage) (ok bool) {
	store.mx.Lock()
	defer store.mx.Unlock()
	rec, ok := store.records[id]
	if ok {
		rec.Data = data
		store.records[id] = rec
	}
	return
}

// DeleteRecord delete record by ID
func DeleteRecord(id string) bool {
	store.mx.Lock()
	defer store.mx.Unlock()
	_, ok := store.records[id]
	if ok {
		delete(store.records, id)
	}
	return ok
}
