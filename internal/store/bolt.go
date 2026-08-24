package store

import (
	"errors"
	"os"

	bolt "go.etcd.io/bbolt"
)

var bucketNames = [][]byte{
	[]byte("profiles"),
	[]byte("batches"),
	[]byte("records"),
	[]byte("issues"),
	[]byte("summaries"),
}

type Store struct {
	db   *bolt.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 0})
	if err != nil {
		return nil, err
	}
	store := &Store{db: db, path: path}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

func (s *Store) Healthy() bool {
	if s == nil || s.db == nil {
		return false
	}
	return s.db.View(func(tx *bolt.Tx) error { return nil }) == nil
}

func ensureOpen(s *Store) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}

func cloneBytes(value []byte) []byte {
	if value == nil {
		return nil
	}
	copyValue := make([]byte, len(value))
	copy(copyValue, value)
	return copyValue
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
