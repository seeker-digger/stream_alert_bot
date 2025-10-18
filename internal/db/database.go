package db

import (
	"encoding/json"
	"errors"
	bolt "go.etcd.io/bbolt"
	l "main.go/internal/logger"
	"slices"
	"strconv"
)

var ErrKeyNotExist = errors.New("key not exist")

type User struct {
	Kick []string
}

func (db *DB) GetUser(id int64) (User, error) {
	user := User{}
	db.mu.RLock() //For safety read!!
	defer db.mu.RUnlock()

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get(itob(id))
		if v == nil {
			return ErrKeyNotExist
		}
		if err := json.Unmarshal(v, &user); err != nil {
			l.Log.Panicf("unsuccessful marshalization: %v", err)
		}
		return nil
	})
	if errors.Is(err, ErrKeyNotExist) {
		return User{}, ErrKeyNotExist
	} else if err != nil {
		return User{}, errors.New("DB view error: " + err.Error())
	}
	l.Log.Debugln("Received user: " + strconv.FormatInt(id, 10))
	return user, nil
}

func (db *DB) SetUser(id int64, user User) error {
	db.mu.Lock() //For safety write!!
	defer db.mu.Unlock()

	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		m, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put(itob(id), m)
	})
	if err != nil {
		return errors.New("DB set error: " + err.Error())
	}
	l.Log.Debugln("New update user: " + strconv.FormatInt(id, 10))
	return nil
}

func (db *DB) RemoveUser(id int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.Delete(itob(id))
	})
	if err != nil {
		return errors.New("DB delete error: " + err.Error())
	}
	l.Log.Debugln("Removed user: " + strconv.FormatInt(id, 10))
	return nil
}

func (db *DB) GetAllIdsByValueKick(val string) ([]int64, error) {
	var ids []int64
	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user User
			if err := json.Unmarshal(v, &user); err != nil {
				continue
			}
			if slices.Contains(user.Kick, val) {
				ids = append(ids, btoi(k))
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("DB view error: " + err.Error())
	}
	return ids, nil
}

func (db *DB) GetCountOfUsers() (int, error) {
	var count = 0
	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, errors.New("DB view error: " + err.Error())
	}
	return count, nil
}

func (db *DB) GetAllUsers() (map[int64]User, error) {
	users := make(map[int64]User)
	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user User
			if err := json.Unmarshal(v, &user); err != nil {
				continue
			}
			users[btoi(k)] = user
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("DB view error: " + err.Error())
	}
	return users, nil
}

func (db *DB) GetAllUniqueValues() ([]string, error) {
	var slugs []string
	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		c := b.Cursor()
		var u User
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &u)
			if err != nil {
				continue
			}
			for _, slug := range u.Kick {
				if !slices.Contains(slugs, slug) {
					slugs = append(slugs, slug)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("DB view error: " + err.Error())
	}
	return slugs, nil
}
