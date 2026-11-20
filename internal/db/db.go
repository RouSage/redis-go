package db

import "sync"

type DB struct {
	store map[string]string
	mu    sync.Mutex
}

func New() *DB {
	return &DB{
		store: make(map[string]string),
	}
}

func (db *DB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.store[key] = value
}

func (db *DB) Get(key string) string {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.store[key]
}
