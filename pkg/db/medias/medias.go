package medias

import (
	"fmt"

	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
	"github.com/sirupsen/logrus"
)

type Store struct {
	store db.StoreBase
}

func CreateStore(database db.Client) (*Store, error) {
	store, err := database.CreateStore(func() db.Entity {
		return new(Media)
	})

	if err != nil {
		return nil, err
	}

	return &Store{store}, nil
}

func (s *Store) Transactional(tx db.Transaction) *Store {
	return &Store{s.store.Transactional(tx)}
}

func (s *Store) List() ([]Media, error) {
	var medias = []Media{}
	return medias, s.store.List(&medias)
}

func (s *Store) Get(id int) (*Media, error) {
	var result = new(Media)
	return result, s.store.Get(id, &result)
}

func (s *Store) Insert(m *Media) (err error) {
	logrus.Debugf("Before ensure")
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		if err = store.Insert(m); err != nil {
			return err
		}
		logrus.Debugf("After insert %v", m)

		// Set new media name if not given
		if m.Name == "" {
			m.Name = fmt.Sprintf("Media %d", m.ID)
			if err := store.Update(m); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) Update(m *Media) error {
	return s.store.Update(m)
}

func (s *Store) Delete(id int) error {
	return s.store.Delete(id)
}

func (s *Store) UpsertAll(medias []Media) error {
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		for _, m := range medias {
			if err := store.Upsert(&m); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) Exists(id int) (bool, error) {
	return s.store.Exists(id)
}
