package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
)

type store struct {
	table         string
	idType        string
	entityFactory func() db.Entity
	pg            *sql.DB
	tx            *sql.Tx
}

type storeClient interface {
	Exec(query string, params ...interface{}) (sql.Result, error)
	Query(query string, params ...interface{}) (*sql.Rows, error)
	QueryRow(query string, params ...interface{}) *sql.Row
}

var _ storeClient = (*sql.DB)(nil)
var _ storeClient = (*sql.Tx)(nil)

func (s *store) client() storeClient {
	if s.tx != nil {
		return s.tx
	}
	return s.pg
}

var _ db.StoreBase = new(store)

func (s *store) WithTransaction(tx db.Transaction) db.StoreBase {
	if tx == nil {
		return s
	}

	pgTx, ok := tx.(*sql.Tx)
	if !ok {
		panic("Postgres driver received non postgres transaction (*sql.Tx required)")
	}

	return &store{s.table, s.idType, s.entityFactory, s.pg, pgTx}
}

func (s *store) IsWithTransaction() bool {
	return s.tx != nil
}

func (s *store) selectQuery(query string) string {
	return fmt.Sprintf(`SELECT id, data FROM "%s" %s`, s.table, query)
}

func (s *store) Get(id interface{}, result interface{}) error {
	row := s.client().QueryRow(s.selectQuery("WHERE id = $1"), id)

	if err := s.unmarshallRow(row, result); err != nil {
		return err
	}

	return nil
}

func (s *store) Exists(id interface{}) (bool, error) {
	entity := s.entityFactory()

	if err := s.Get(id, &entity); err != nil {
		return false, err
	}

	return entity != nil, nil
}

func (s *store) List(result interface{}) error {
	return s.Find(nil, result)
}

func (s *store) Find(filters map[string]interface{}, result interface{}) error {
	resultVal := reflect.ValueOf(result)
	// FIXME remove this test
	if resultVal.Kind() != reflect.Ptr || resultVal.Elem().Kind() != reflect.Slice {
		panic("List result should be a slice pointer")
	}
	sliceVal := resultVal.Elem()

	if filters == nil {
		filters = map[string]interface{}{}
	}

	jsonFilters, err := json.Marshal(filters)
	if err != nil {
		return nil
	}

	rows, err := s.client().Query(s.selectQuery("WHERE data @> $1::jsonb"), jsonFilters)
	if err != nil {
		return err
	}

	for rows.Next() {
		entity := s.entityFactory()
		if err := s.unmarshallRow(rows, &entity); err != nil {
			return err
		}
		sliceVal = reflect.Append(sliceVal, reflect.ValueOf(entity).Elem())
	}

	resultVal.Elem().Set(sliceVal.Slice(0, sliceVal.Len()))
	return nil
}

func (s *store) Insert(item db.Entity) error {
	id, data := prepare(item)

	if db.ShouldAutoIncrement(item) {
		row := s.client().QueryRow(fmt.Sprintf(`INSERT INTO "%s" (data) VALUES ($1) RETURNING id`, s.table), data)
		if err := row.Scan(&id); err != nil {
			return err
		}
		item.SetID(id)
	} else {
		if _, err := s.client().Exec(fmt.Sprintf(`INSERT INTO "%s" (id, data) VALUES ($1, $2)`, s.table), id, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) Update(item db.Entity) error {
	id, data := prepare(item)

	result, err := s.client().Exec(fmt.Sprintf(`UPDATE "%s" SET data = $1 WHERE id = $2`, s.table), data, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return db.EntityNotFound
	}

	return nil
}

func (s *store) Upsert(item db.Entity) error {
	id, data := prepare(item)

	if db.ShouldAutoIncrement(item) {
		row := s.client().QueryRow(fmt.Sprintf(`INSERT INTO "%s" (data) VALUES ($1) RETURNING id`, s.table), data)
		if err := row.Scan(&id); err != nil {
			return err
		}
		item.SetID(id)
	} else {
		if _, err := s.client().Exec(fmt.Sprintf(`
				INSERT INTO "%s" (id, data) VALUES ($1, $2)
				ON CONFLICT (id) DO UPDATE SET data = $2
			`, s.table), id, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) Delete(id interface{}) error {
	_, err := s.client().Exec(fmt.Sprintf(`DELETE FROM "%s" WHERE id = $1`, s.table), id)
	return err
}

func (s *store) DeleteAll() error {
	_, err := s.client().Exec(fmt.Sprintf(`DELETE FROM "%s"`, s.table))
	return err
}

func (s *store) unmarshallRow(row scanable, result interface{}) error {
	if s.idType == "serial" {
		return unmarshallRow(row, 0, result)
	}
	return unmarshallRow(row, "", result)
}
