package store

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

const (
	backupFile     = ".store.gob"
	backupFileTemp = ".store.go.TMP"
)

var (
	ticker      *time.Ticker
	opsCountCfg int
)

func runBackupByTicker(secsInterval int) {
	if secsInterval > 0 {
		d := time.Second * time.Duration(secsInterval)
		ticker = time.NewTicker(d)
		go func() {
			for range ticker.C {
				go backupStore()
			}
		}()
	}
}

func runBackupByOpsCounter() {
	if opsCountCfg > 0 && store.opsCount >= opsCountCfg {
		store.mx.Lock()
		store.opsCount = 0
		store.mx.Unlock()
		go backupStore()
	}
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
	if store.backupInProcess == false {
		updateBackupStateKey(true)
	} else {
		return
	}
	fmt.Println("Saving store to disk ...")
	f, err := os.Create(backupFileTemp)
	if err == nil {
		defer f.Close()
		enc := gob.NewEncoder(f)
		enc.Encode(store.records)
		err = os.Rename(backupFileTemp, backupFile)
		if err != nil {
			fmt.Printf("Unable to save %s file", backupFile)
		} else {
			fmt.Println("Store saved to disk")
		}
		updateBackupStateKey(false)
	}
	return
}

func updateBackupStateKey(newState bool) {
	store.mx.Lock()
	store.backupInProcess = newState
	store.mx.Unlock()
}
