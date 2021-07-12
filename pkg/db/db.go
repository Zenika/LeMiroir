package db

import (
	"github.com/allez-chauffe/marcel/pkg/config"
	"github.com/allez-chauffe/marcel/pkg/db/clients"
	"github.com/allez-chauffe/marcel/pkg/db/driver/driver"
	_ "github.com/allez-chauffe/marcel/pkg/db/driver/register" // Register drivers
	"github.com/allez-chauffe/marcel/pkg/db/medias"
	"github.com/allez-chauffe/marcel/pkg/db/plugins"
	"github.com/allez-chauffe/marcel/pkg/db/users"
)

type DB struct {
	client  driver.Client
	clients *clients.Store
	medias  *medias.Store
	plugins *plugins.Store
	users   *users.Store
}

func Open() (database *DB, err error) {
	database = new(DB)

	database.client, err = driver.Get(config.Default().API().DB().Driver()).Open()
	if err != nil {
		return nil, err
	}

	if database.clients, err = clients.NewStore(database.client); err != nil {
		return nil, err
	}
	if database.medias, err = medias.CreateStore(database.client); err != nil {
		return nil, err
	}
	if database.plugins, err = plugins.CreateStore(database.client); err != nil {
		return nil, err
	}
	if database.users, err = users.CreateStore(database.client); err != nil {
		return nil, err
	}

	return
}

func (database *DB) Close() error {
	return database.client.Close()
}

func (database *DB) Begin() (*Tx, error) {
	tx, err := database.client.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

// FIXME
// func Transactional(task func(*Tx) error) (err error) {
// 	return driver.Transactional(func(tx driver.Transaction) error {
// 		return task(&Tx{tx})
// 	})
// }
