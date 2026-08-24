package store

import (
	"errors"
	"fmt"

	"example.com/childfitness/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveProfile(profile domain.ChildProfile) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	if err := domain.ValidateProfile(profile); err != nil {
		return err
	}
	data, err := encode(profile)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("profiles")).Put(key(profile.Key()), data)
	})
}

func (s *Store) GetProfile(schoolID, classID, childID string) (domain.ChildProfile, error) {
	if err := ensureOpen(s); err != nil {
		return domain.ChildProfile{}, err
	}
	var profile domain.ChildProfile
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("profiles")).Get(key(domain.ChildProfile{SchoolID: schoolID, ClassID: classID, ID: childID}.Key()))
		if data == nil {
			return errors.New("profile not found")
		}
		return decode(cloneBytes(data), &profile)
	})
	return profile, err
}

func (s *Store) ListProfiles(schoolID, classID string) ([]domain.ChildProfile, error) {
	if err := ensureOpen(s); err != nil {
		return nil, err
	}
	profiles := make([]domain.ChildProfile, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("profiles")).ForEach(func(_, data []byte) error {
			var profile domain.ChildProfile
			if err := decode(data, &profile); err != nil {
				return err
			}
			if profile.SchoolID == schoolID && profile.ClassID == classID {
				profiles = append(profiles, profile)
			}
			return nil
		})
	})
	return profiles, err
}

func (s *Store) SaveBatch(batch domain.FitnessBatch) error {
	if err := ensureOpen(s); err != nil {
		return err
	}
	if err := domain.ValidateBatch(batch); err != nil {
		return err
	}
	batchData, err := encode(batch)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte("batches")).Put(key(batch.Key()), batchData); err != nil {
			return err
		}
		for _, record := range batch.Records {
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := tx.Bucket([]byte("records")).Put(key(record.Key()+":"+record.ID), data); err != nil {
				return err
			}
		}
		for index, issue := range batch.Issues {
			data, err := encode(issue)
			if err != nil {
				return err
			}
			issueKey := fmt.Sprintf("%s:%d:%d", batch.ID, issue.Row, index)
			if err := tx.Bucket([]byte("issues")).Put(key(issueKey), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) GetBatch(batchID string) (domain.FitnessBatch, error) {
	if err := ensureOpen(s); err != nil {
		return domain.FitnessBatch{}, err
	}
	var batch domain.FitnessBatch
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("batches")).Get(key(batchID))
		if data == nil {
			return errors.New("batch not found")
		}
		return decode(cloneBytes(data), &batch)
	})
	return batch, err
}

func (s *Store) ListBatches(schoolID, classID string) ([]domain.FitnessBatch, error) {
	if err := ensureOpen(s); err != nil {
		return nil, err
	}
	batches := make([]domain.FitnessBatch, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("batches")).ForEach(func(_, data []byte) error {
			var batch domain.FitnessBatch
			if err := decode(data, &batch); err != nil {
				return err
			}
			if batch.SchoolID == schoolID && batch.ClassID == classID {
				batches = append(batches, batch)
			}
			return nil
		})
	})
	return batches, err
}
