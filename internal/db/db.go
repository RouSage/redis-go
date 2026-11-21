package db

import (
	"errors"
	"sync"
	"time"
)

var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

type value interface {
	isValue()
}

type stringValue string
type listValue []string

func (stringValue) isValue() {}
func (listValue) isValue()   {}

type data struct {
	value     value
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

	d := data{value: stringValue(value)}
	if ttl > 0 {
		d.expiresAt = time.Now().Add(ttl)
	}
	db.store[key] = d
}

func (db *DB) Get(key string) (string, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	d, ok := db.lookup(key)
	if !ok {
		return "", ok, nil
	}

	s, ok := d.value.(stringValue)
	if !ok {
		return "", ok, ErrWrongType
	}

	return string(s), ok, nil
}

func (db *DB) RPush(key string, values ...string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// TODO: keep TTL when EXPIRE exists
	var list listValue
	if d, ok := db.lookup(key); ok {
		l, ok := d.value.(listValue)
		if !ok {
			return 0, ErrWrongType
		}

		list = l
	}

	list = append(list, values...)
	db.store[key] = data{value: list}

	return len(list), nil
}

func (db *DB) lookup(key string) (data, bool) {
	d, ok := db.store[key]
	if !ok {
		return data{}, ok
	}

	// lazy expiry on read
	if d.expired(time.Now()) {
		db.delete(key)
		return data{}, false
	}

	return d, ok
}

func (db *DB) delete(key string) {
	delete(db.store, key)
}
