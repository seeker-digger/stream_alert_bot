package db

import (
	"errors"
	"log"
	"testing"
)

func TestDB_SetUser(t *testing.T) {
	db := Init()
	defer db.db.Close()
	var list = []string{"a", "jesusavgn", "c"}
	err := db.SetUser(1234, User{Kick: list})
	if err != nil {
		t.Error(err)
	}
}

func TestDB_GetUser(t *testing.T) {
	db := Init()
	defer db.db.Close()
	u, err := db.GetUser(1195310723)
	if errors.Is(err, ErrKeyNotExist) {
		t.Log("User not exist")
	} else if err != nil {
		t.Error(err)
	}

	t.Log(u)
}

func TestDB_RemoveUser(t *testing.T) {
	db := Init()
	defer db.db.Close()
	err := db.RemoveUser(1234)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_GetAllIdsByValueKick(t *testing.T) {
	db := Init()
	defer db.db.Close()

	u, err := db.GetUser(1234)
	if err != nil {
		t.Error(err)
	}
	u.Kick = append(u.Kick, "jesusavgn")
	err = db.SetUser(1234, u)
	if err != nil {
		t.Error(err)
	}
	ids, err := db.GetAllIdsByValueKick("jesusavgn")
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
}

func TestDB_GetCountOfUsers(t *testing.T) {
	db := Init()
	defer db.db.Close()

	c, err := db.GetCountOfUsers()
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}

func TestDB_GetAllUsers(t *testing.T) {
	db := Init()
	defer db.db.Close()
	u, err := db.GetAllUsers()
	if err != nil {
		t.Error(err)
	}
	for k, v := range u {
		t.Log(k, v)
	}
}

func TestDB_GetAllUniqueValues(t *testing.T) {
	db := Init()
	defer db.db.Close()
	u, err := db.GetAllUniqueValues()
	if err != nil {
		t.Error(err)
	}
	log.Println(u)
}
