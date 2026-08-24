package store

import (
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) RecordsForClass(schoolID, classID string) ([]domain.MeasurementRecord, error) {
	profiles, err := s.ListProfiles(schoolID, classID)
	if err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(profiles))
	for _, profile := range profiles {
		allowed[profile.ID] = true
	}
	records := make([]domain.MeasurementRecord, 0)
	err = s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, data []byte) error {
			var record domain.MeasurementRecord
			if err := decode(data, &record); err != nil {
				return err
			}
			if allowed[record.ChildID] {
				record.Status = strings.TrimSpace(record.Status)
				records = append(records, record)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	return records, nil
}

func (s *Store) SaveSummary(summary domain.ClassSummary) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	data, err := encode(summary)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("summaries")).Put(key(summary.SchoolID+":"+summary.ClassID), data)
	})
}

func (s *Store) GetSummary(schoolID, classID string) (domain.ClassSummary, error) {
	if err := ensureOpen(s); err != nil {
		return domain.ClassSummary{}, err
	}
	var summary domain.ClassSummary
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("summaries")).Get(key(schoolID + ":" + classID))
		if data == nil {
			return nil
		}
		return decode(cloneBytes(data), &summary)
	})
	return summary, err
}

func (s *Store) Count(bucketName string) (int, error) {
	if err := ensureOpen(s); err != nil {
		return 0, err
	}
	count := 0
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}
		count = bucket.Stats().KeyN
		return nil
	})
	return count, err
}
