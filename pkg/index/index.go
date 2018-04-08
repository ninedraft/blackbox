package index

import (
	"errors"
)

type KV interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
}

var (
	ErrValueNotFound = errors.New("value is not found")
)

type mapKV map[string]string

func (kv mapKV) Set(k, v []byte) error {
	kv[string(k)] = string(v)
	return nil
}

func (kv mapKV) Get(k []byte) ([]byte, error) {
	val, isFound := kv[string(k)]
	if !isFound {
		return nil, ErrValueNotFound
	}
	return []byte(val), nil
}

type Indexer struct {
	kv KV
}
