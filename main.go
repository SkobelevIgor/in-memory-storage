package main

import (
	"flag"
	"in-memory-storage/api"
	"in-memory-storage/store"
	"log"
	"net/http"
)

func main() {
	var storeCfg store.Config
	flag.IntVar(&storeCfg.BackupTimeSecs, "backup-interval", 0, "Backup store interval (secs)")
	flag.IntVar(&storeCfg.BackupOps, "backup-ops-count", 0, "Backup each insert / update ops")
	flag.BoolVar(&storeCfg.GracefulBackup, "graceful-backup", false, "Backup store on shutdown")
	flag.BoolVar(&storeCfg.SkipBackupLoad, "skip-backup-load", false, "Load store from backup")
	flag.Parse()
	err := store.InitStore(storeCfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", api.RequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
