package bolt

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/bolthold"

	"github.com/allez-chauffe/marcel/pkg/db/driver/driver"
)

type client struct {
	bh *bh.Store
}

func (c *client) Store(options ...driver.StoreOption) (driver.Store, error) {
	return driver.NewStore(func(config *driver.StoreConfig) driver.StoreBase {
		return &store{c, config, nil}
	}, options...)
}

func (c *client) Begin() (driver.Transaction, error) {
	return c.bh.Bolt().Begin(true)
}

func (c *client) Close() error {
	log.Info("Closing bolt database...")

	err := c.bh.Close()
	if err != nil {
		return fmt.Errorf("Error while closing bolt database: %w", err)
	}

	log.Info("bolt database closed")

	return nil
}
