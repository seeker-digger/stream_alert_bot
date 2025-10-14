package db

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"main.go/internal/config"
	"sync"
)

type DB struct {
	db *bolt.DB
	mu sync.RWMutex
}

func Init() DB {
	db, err := bolt.Open(config.GetDataPath("app.db"), 0660, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		return err
	}); err != nil {
		log.Fatal(err)
	}
	return DB{db: db}
}
