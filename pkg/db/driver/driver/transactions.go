package driver

import (
	log "github.com/sirupsen/logrus"
)

type Transaction interface {
	Commit() error
	Rollback() error
}

func Transactional(client Client, task func(Transaction) error) (err error) {
	tx, err := client.Begin()
	if err != nil {
		return err
	}
	log.Debugf("begin transaction")

	defer finishTransaction(tx, &err)

	return task(tx)
}

func finishTransaction(tx Transaction, err *error) {
	r := recover()

	if *err == nil && r == nil {
		*err = tx.Commit()
		log.Debugf("commit transaction")
		return
	}

	if err := tx.Rollback(); err != nil {
		log.Errorf("Error while rolling back: %s", err)
	}
	log.Debugf("rollback transaction")

	if r != nil {
		panic(r)
	}
}
