package clients

import "github.com/allez-chauffe/marcel/pkg/db/driver/driver"

type Store struct {
	store driver.Store
}

func NewStore(client driver.Client) (*Store, error) {
	store, err := client.Store(driver.WithType(new(Client)))
	if err != nil {
		return nil, err
	}

	return &Store{store}, nil
}

// FIXME
// func (s *Store) Transactional(tx db.Transaction) *Store {
// 	return &Store{s.store.Transactional(tx)}
// }

func (s *Store) Get(id string) (*Client, error) {
	e, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return e.(*Client), nil
}

func (s *Store) Exists(id string) (bool, error) {
	return s.store.Exists(id)
}

func (s *Store) List() ([]Client, error) {
	l, err := s.store.List()
	if err != nil {
		return nil, err
	}
	return l.([]Client), nil
}

func (s *Store) Insert(c *Client) error {
	return s.store.Insert(c)
}

func (s *Store) Update(c *Client) error {
	return s.store.Update(c)
}

func (s *Store) Delete(id string) error {
	return s.store.Delete(id)
}

func (s *Store) DeleteAll() error {
	return s.store.DeleteAll()
}
