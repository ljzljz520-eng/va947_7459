package store

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type DatabaseStats struct {
	Path       string         `json:"path"`
	Healthy    bool           `json:"healthy"`
	EntityRows map[string]int `json:"entity_rows"`
}

func (s *Store) Stats() (DatabaseStats, error) {
	if err := ensureOpen(s); err != nil {
		return DatabaseStats{}, err
	}
	stats := DatabaseStats{Path: s.path, Healthy: s.Healthy(), EntityRows: make(map[string]int)}
	err := s.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range []string{"profiles", "batches", "records", "issues", "summaries"} {
			stats.EntityRows[bucket] = tx.Bucket([]byte(bucket)).Stats().KeyN
		}
		return nil
	})
	return stats, err
}

func (s *Store) ExportJSON() ([]byte, error) {
	stats, err := s.Stats()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(stats, "", "  ")
}

func (s *Store) VerifyBuckets() error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	return s.db.View(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if tx.Bucket(name) == nil {
				return errors.New("missing bucket " + string(name))
			}
		}
		return nil
	})
}

func (s *Store) Keys(bucketName string) ([]string, error) {
	if err := ensureOpen(s); err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("unknown bucket " + bucketName)
		}
		return bucket.ForEach(func(key, _ []byte) error {
			keys = append(keys, string(key))
			return nil
		})
	})
	sort.Strings(keys)
	return keys, err
}

func (s *Store) DeleteBatch(batchID string) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte("batches")).Delete(key(batchID)); err != nil {
			return err
		}
		for _, bucketName := range []string{"records", "issues"} {
			bucket := tx.Bucket([]byte(bucketName))
			toDelete := make([][]byte, 0)
			if err := bucket.ForEach(func(itemKey, _ []byte) error {
				if strings.HasPrefix(string(itemKey), batchID+":") {
					toDelete = append(toDelete, cloneBytes(itemKey))
				}
				return nil
			}); err != nil {
				return err
			}
			for _, itemKey := range toDelete {
				if err := bucket.Delete(itemKey); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Store) DeleteClassSummary(schoolID, classID string) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("summaries")).Delete(key(schoolID + ":" + classID))
	})
}
