package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// Import of postgres driver
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/allez-chauffe/marcel/pkg/config"
	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
)

type driver struct{}

// Driver is the postgres driver implementation.
var Driver db.Driver = new(driver)

func (d *driver) Open() (db.Client, error) {
	pgConf := config.Default().API().DB().Postgres()

	log.Infof("Connecting to postgres database (%s:%s/%s) ...", pgConf.Host(), pgConf.Port(), pgConf.DBName())

	pg, err := sql.Open("postgres", getConnectionString(""))
	if err != nil {
		return nil, err
	}

	if err := pg.Ping(); err != nil {
		if strings.Contains(err.Error(), "database") && strings.Contains(err.Error(), "does not exist") {
			pg.Close()
			var creationError error
			if pg, creationError = createDatabase(); creationError != nil {
				return nil, fmt.Errorf("%w (failed to create it: %s", err, creationError)
			}
		} else {
			return nil, err
		}
	}

	log.Info("Postgres database connected")

	return &client{pg}, nil
}

type client struct {
	pg *sql.DB
}

func (c *client) Begin() (db.Transaction, error) {
	return c.pg.Begin()
}

func (c *client) Store(newEntity func() db.Entity) (db.Store, error) {
	entity := newEntity()
	table := toSnakeCase(reflect.TypeOf(entity).Elem().Name())
	idType := reflect.TypeOf(entity.GetID()).Name()

	var postgresIDType string
	if strings.HasPrefix(idType, "int") {
		postgresIDType = "serial"
	} else if idType == "string" {
		postgresIDType = "text"
	} else {
		return nil, errors.New("Postgres database driver only supports uuid and int ID types")
	}

	log.Infof("Ensure table %s exists", table)
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
			"id" %s NOT NULL PRIMARY KEY,
			"data" jsonb NOT NULL
		);
	`, table, postgresIDType)

	if _, err := c.pg.Exec(query); err != nil {
		return nil, fmt.Errorf("Error while creating table '%s': %w \nQuery: %s", table, err, query)
	}

	return &store{table, postgresIDType, newEntity, c.pg, nil}, nil
}

func (c *client) Close() error {
	return c.pg.Close()
}
