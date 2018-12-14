package store

import (
	"time"
)

type Record struct {
	Id        int64
	Data      string
	CreatedAt *time.Time
}

func GetRecords() (records []Record, err Error) {

}

func GetRecord(id int64) (record Record, err Error) {

}

func SaveRecord() (id int64, err Error) {

}

func UpdateRecord(id int64, data string) {

}

func DeleteRecord(id int64) {

}
