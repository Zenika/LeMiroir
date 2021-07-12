package bolt

import (
	bh "github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"

	"github.com/allez-chauffe/marcel/pkg/db/driver/driver"
)

type store struct {
	*client
	config *driver.StoreConfig
	tx     *bolt.Tx
}

var _ driver.StoreBase = new(store)

func (s *store) Get(id interface{}, e interface{}) (err error) {
	if s.HasTransaction() {
		err = s.bh.TxGet(s.tx, id, e)
	} else {
		err = s.bh.Get(id, e)
	}

	if err == bh.ErrNotFound {
		err = driver.ErrNotFound
	}

	return
}

func (s *store) Insert(id interface{}, e interface{}) error {
	if s.config.AutoIncrement() {
		id = bh.NextSequence()
	}

	if s.HasTransaction() {
		return s.bh.TxInsert(s.tx, id, e)
	}
	return s.bh.Insert(id, e)
}

func (s *store) Update(id interface{}, e interface{}) (err error) {
	if s.HasTransaction() {
		err = s.bh.TxUpdate(s.tx, id, e)
	} else {
		err = s.bh.Update(id, e)
	}

	if err == bh.ErrNotFound {
		err = driver.ErrNotFound
	}

	return
}

func (s *store) Delete(id interface{}) error {
	if s.HasTransaction() {
		return s.bh.TxDelete(s.tx, id, s.config.New())
	}
	return s.bh.Delete(id, s.config.New())
}

func (s *store) WithTransaction(tx driver.Transaction) driver.StoreBase {
	if tx == nil { // FIXME is this useful?
		return s
	}

	if s.HasTransaction() {
		return s // FIXME check if same tx?
	}

	boltTx, ok := tx.(*bolt.Tx)
	if !ok {
		panic("Bolt driver received non bolt transaction (*bbolt.Tx required)")
	}

	return &store{s.client, s.config, boltTx}
}

func (s *store) HasTransaction() bool {
	return s.tx != nil
}

var _ driver.StoreDeleteAll = new(store)

func (s *store) DeleteAll() error {
	if s.HasTransaction() {
		return s.bh.TxDeleteMatching(s.tx, s.config.New(), nil)
	}
	return s.bh.DeleteMatching(s.config.New(), nil)
}

var _ driver.StoreUpsert = new(store)

func (s *store) Upsert(id, e interface{}) error {
	if id == uint64(0) && s.config.AutoIncrement() {
		id = bh.NextSequence()
	}

	if s.HasTransaction() {
		return s.bh.TxUpsert(s.tx, id, e)
	}
	return s.bh.Upsert(id, e)
}

var _ driver.StoreFind = new(store)

func (s *store) Find(filters map[string]interface{}, result interface{}) error {
	query := queryFromFilters(filters)
	if s.HasTransaction() {
		return s.bh.TxFind(s.tx, result, query)
	}
	return s.bh.Find(result, query)
}
