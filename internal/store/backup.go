package store

import (
	"encoding/json"
	"errors"
	"sort"

	"example.com/childfitness/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type Backup struct {
	Profiles  []domain.ChildProfile `json:"profiles"`
	Batches   []domain.FitnessBatch `json:"batches"`
	Summaries []domain.ClassSummary `json:"summaries"`
}

func (s *Store) CreateBackup() (Backup, error) {
	if err := ensureOpen(s); err != nil {
		return Backup{}, err
	}
	backup := Backup{Profiles: make([]domain.ChildProfile, 0), Batches: make([]domain.FitnessBatch, 0), Summaries: make([]domain.ClassSummary, 0)}
	err := s.db.View(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte("profiles")).ForEach(func(_, data []byte) error {
			var profile domain.ChildProfile
			if err := decode(data, &profile); err != nil {
				return err
			}
			backup.Profiles = append(backup.Profiles, profile)
			return nil
		}); err != nil {
			return err
		}
		if err := tx.Bucket([]byte("batches")).ForEach(func(_, data []byte) error {
			var batch domain.FitnessBatch
			if err := decode(data, &batch); err != nil {
				return err
			}
			backup.Batches = append(backup.Batches, batch)
			return nil
		}); err != nil {
			return err
		}
		return tx.Bucket([]byte("summaries")).ForEach(func(_, data []byte) error {
			var summary domain.ClassSummary
			if err := decode(data, &summary); err != nil {
				return err
			}
			backup.Summaries = append(backup.Summaries, summary)
			return nil
		})
	})
	sort.Slice(backup.Profiles, func(i, j int) bool { return backup.Profiles[i].Key() < backup.Profiles[j].Key() })
	sort.Slice(backup.Batches, func(i, j int) bool { return backup.Batches[i].ID < backup.Batches[j].ID })
	return backup, err
}

func (s *Store) BackupJSON() ([]byte, error) {
	backup, err := s.CreateBackup()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(backup, "", "  ")
}

func (s *Store) RestoreBackup(backup Backup) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	if len(backup.Profiles) == 0 && len(backup.Batches) == 0 && len(backup.Summaries) == 0 {
		return errors.New("backup is empty")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, profile := range backup.Profiles {
			if err := domain.ValidateProfile(profile); err != nil {
				return err
			}
			data, err := encode(profile)
			if err != nil {
				return err
			}
			if err := tx.Bucket([]byte("profiles")).Put(key(profile.Key()), data); err != nil {
				return err
			}
		}
		for _, batch := range backup.Batches {
			if err := domain.ValidateBatch(batch); err != nil {
				return err
			}
			data, err := encode(batch)
			if err != nil {
				return err
			}
			if err := tx.Bucket([]byte("batches")).Put(key(batch.ID), data); err != nil {
				return err
			}
		}
		for _, summary := range backup.Summaries {
			data, err := encode(summary)
			if err != nil {
				return err
			}
			if err := tx.Bucket([]byte("summaries")).Put(key(summary.SchoolID+":"+summary.ClassID), data); err != nil {
				return err
			}
		}
		return nil
	})
}
