package db

import (
	"sync"
	"time"
)

type data struct {
	value     string
	expiresAt time.Time
}

func (d data) expired(now time.Time) bool {
	return !d.expiresAt.IsZero() && now.After(d.expiresAt)
}

type DB struct {
	store map[string]data
	mu    sync.Mutex
}

func New() *DB {
	return &DB{
		store: make(map[string]data),
	}
}

func (db *DB) Set(key, value string, ttl time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()

	d := data{value: value}
	if ttl > 0 {
		d.expiresAt = time.Now().Add(ttl)
	}
	db.store[key] = d
}

func (db *DB) Get(key string) (string, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	d, ok := db.store[key]
	if !ok {
		return "", ok
	}

	// lazy expiry on read
	if d.expired(time.Now()) {
		db.delete(key)
		return "", false
	}

	return d.value, ok
}

func (db *DB) delete(key string) {
	delete(db.store, key)
}
