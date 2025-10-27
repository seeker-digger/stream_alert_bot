package db

import (
	bolt "go.etcd.io/bbolt"
	"main.go/internal/config"
	l "main.go/internal/logger"
	"sync"
)

type DB struct {
	db *bolt.DB
	mu sync.RWMutex
}

func Init() DB {
	db, err := bolt.Open(config.GetDataPath("app.db"), 0666, nil)
	if err != nil {
		l.Log.Fatal(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		return err
	}); err != nil {
		l.Log.Fatal(err)
	}
	return DB{db: db}
}
