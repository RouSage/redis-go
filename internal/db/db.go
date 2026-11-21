package db

import (
	"errors"
	"slices"
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
	list, err := db.lookupList(key)
	if err != nil {
		return 0, err
	}

	list = append(list, values...)
	db.store[key] = data{value: list}

	return len(list), nil
}

func (db *DB) LPush(key string, values ...string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// TODO: keep TTL when EXPIRE exists
	list, err := db.lookupList(key)
	if err != nil {
		return 0, err
	}

	for _, value := range values {
		list = append([]string{value}, list...)
	}

	db.store[key] = data{value: list}

	return len(list), nil
}

func (db *DB) LRange(key string, start, stop int) ([]string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	list, err := db.lookupList(key)
	if err != nil {
		return nil, err
	}

	if start < 0 {
		start = max(0, len(list)+start)
	}
	if stop < 0 {
		stop = max(0, len(list)+stop)
	}

	if start > stop || start >= len(list) {
		return []string{}, nil
	}
	if stop >= len(list) {
		stop = len(list) - 1
	}

	return slices.Clone(list[start : stop+1]), nil
}

func (db *DB) LLen(key string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	list, err := db.lookupList(key)
	if err != nil {
		return 0, err
	}

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

func (db *DB) lookupList(key string) (listValue, error) {
	d, ok := db.lookup(key)
	if !ok {
		return nil, nil
	}

	list, ok := d.value.(listValue)
	if !ok {
		return nil, ErrWrongType
	}

	return list, nil
}

func (db *DB) delete(key string) {
	delete(db.store, key)
}
