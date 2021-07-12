package plugins

import (
	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
)

type Store struct {
	store db.StoreBase
}

func CreateStore(database db.Client) (*Store, error) {
	store, err := database.CreateStore(func() db.Entity {
		return new(Plugin)
	})

	if err != nil {
		return nil, err
	}

	return &Store{store}, nil
}

func (s *Store) Transactional(tx db.Transaction) *Store {
	return &Store{s.store.Transactional(tx)}
}

func (s *Store) List() ([]Plugin, error) {
	var plugins = []Plugin{}
	return plugins, s.store.List(&plugins)
}

func (s *Store) Get(eltName string) (*Plugin, error) {
	var p = new(Plugin)
	return p, s.store.Get(eltName, &p)
}

func (s *Store) Exists(eltName string) (bool, error) {
	return s.store.Exists(eltName)
}

func (s *Store) Insert(p *Plugin) error {
	return s.store.Insert(p)
}

func (s *Store) Update(p *Plugin) error {
	return s.store.Update(p)
}

func (s *Store) UpsertAll(plugins []Plugin) error {
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		for _, p := range plugins {
			if err := store.Upsert(&p); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) Delete(eltName string) error {
	return s.store.Delete(eltName)
}
