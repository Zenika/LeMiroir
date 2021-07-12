package users

import (
	"time"

	rand "github.com/Pallinder/go-randomdata"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
)

type Store struct {
	store db.StoreBase
}

func CreateStore(database db.Client) (*Store, error) {
	store, err := database.CreateStore(func() db.Entity {
		return new(User)
	})

	if err != nil {
		return nil, err
	}

	return &Store{store}, nil
}

func (s *Store) Transactional(tx db.Transaction) *Store {
	return &Store{s.store.Transactional(tx)}
}

func (s *Store) EnsureOneUser() error {
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		users := &[]User{}

		if err := store.List(users); err != nil || len(*users) != 0 {
			return err
		}

		log.Info("No users in database, creating admin...")

		u := &User{
			ID:          uuid.NewV4().String(),
			DisplayName: "Admin",
			Login:       "admin",
			Role:        "admin",
			CreatedAt:   time.Now(),
		}

		password := rand.RandStringRunes(10)
		if err := u.SetPassword(password); err != nil {
			return err
		}

		if err := store.Insert(u); err != nil {
			return err
		}

		log.Infof("User admin created with password %s", password)

		return nil
	})
}

func (s *Store) Insert(u *User) error {
	return s.store.Insert(u)
}

func (s *Store) List() ([]User, error) {
	var users []User
	return users, s.store.List(&users)
}

func (s *Store) Get(id string) (*User, error) {
	u := new(User)
	return u, s.store.Get(id, &u) // FIXME non
}

func (s *Store) GetByLogin(login string) (*User, error) {
	var users []User

	filters := map[string]interface{}{
		"Login": login,
	}

	if err := s.store.Find(&users, filters); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func (s *Store) Delete(id string) error {
	return s.store.Delete(id)
}

func (s *Store) Disconnect(id string) error {
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		u := &User{}
		if err := store.Get(id, &u); u == nil || err != nil {
			return err
		}

		u.LastDisconnection = time.Now()
		return store.Update(u)
	})
}

func (s *Store) Update(user *User) error {
	return s.store.Update(user)
}

func (s *Store) UpsertAll(users []User) error {
	return db.EnsureTransaction(s.store, func(store db.Store) error {
		for _, u := range users {
			if err := store.Upsert(&u); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) Exists(id string) (bool, error) {
	return s.store.Exists(id)
}
