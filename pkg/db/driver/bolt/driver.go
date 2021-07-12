package bolt

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"

	"github.com/allez-chauffe/marcel/pkg/config"
	"github.com/allez-chauffe/marcel/pkg/db/driver/driver"
)

func init() {
	driver.Register("bolt", new(boltDriver))
}

type boltDriver struct{}

func (d *boltDriver) Open() (driver.Client, error) {
	log.Info("Opening bolt database...")

	var options = *bolt.DefaultOptions
	options.Timeout = 100 * time.Millisecond

	bhs, err := bh.Open(os.ExpandEnv(config.Default().API().DB().Bolt().File()), 0644, &bh.Options{
		Options: &options,
	})
	if err != nil {
		return nil, fmt.Errorf("Error while opening bolt database: %w", err)
	}

	log.Info("bolt database opened")

	return &client{bhs}, nil
}
