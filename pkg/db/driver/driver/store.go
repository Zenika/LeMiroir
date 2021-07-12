package driver

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrNotImplemented = errors.New("not implemented")
)

type Store interface {
	Get(id interface{}) (interface{}, error)

	Exists(id interface{}) (bool, error)

	List() (interface{}, error)

	Find(filters map[string]interface{}) (interface{}, error)

	Insert(e Entity) error

	Update(e Entity) error

	Upsert(e Entity) error

	Delete(id interface{}) error

	DeleteAll() error

	// FIXME Transactional() or not
}

func NewStore(baseFactory func(config *StoreConfig) StoreBase, options ...StoreOption) (Store, error) {
	var config = new(StoreConfig)

	for _, option := range options {
		option(config)
	}

	var base = baseFactory(config)

	return &store{base, config}, nil
}

type StoreOption func(*StoreConfig) error

func WithType(v interface{}) StoreOption {
	return func(c *StoreConfig) error {
		var t = reflect.TypeOf(v)
		if t != nil && t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t == nil || t.Kind() != reflect.Struct {
			return fmt.Errorf("%#v should be of struct{} or &struct{} kind", v)
		}

		// FIXME check pointer to t is AssignableTo Entity

		c.new = func() interface{} {
			return reflect.New(t).Interface()
		}

		var st = reflect.SliceOf(t)
		c.newList = func() interface{} {
			var p = reflect.New(st)
			p.Elem().Set(reflect.MakeSlice(st, 0, 10))
			return p.Interface()
		}

		return nil
	}
}

func WithAutoIncrement() StoreOption {
	return func(c *StoreConfig) error {
		c.autoIncrement = true
		return nil
	}
}

type StoreConfig struct {
	new           func() interface{}
	newList       func() interface{}
	autoIncrement bool
}

func (c *StoreConfig) New() interface{} {
	return c.new()
}

func (c *StoreConfig) NewList() interface{} {
	return c.newList()
}

func (c *StoreConfig) AutoIncrement() bool {
	return c.AutoIncrement()
}

type StoreBase interface {
	Get(id interface{}, e interface{}) error

	Insert(id interface{}, e interface{}) error

	Update(id interface{}, e interface{}) error

	Delete(id interface{}) error

	WithTransaction(tx Transaction) StoreBase

	HasTransaction() bool
}

type StoreExists interface {
	Exists(id interface{}) (bool, error)
}

type StoreFind interface {
	Find(filters map[string]interface{}, result interface{}) error
}

type StoreList interface {
	List(result interface{}) error
}

type StoreUpsert interface {
	Upsert(id interface{}, item interface{}) error
}

type StoreDeleteAll interface {
	DeleteAll() error
}

type Entity interface {
	GetID() interface{}
	SetID(id interface{})
}

type store struct {
	base   StoreBase
	config *StoreConfig
}

func (s *store) Get(id interface{}) (interface{}, error) {
	var e = s.config.new()

	if err := s.base.Get(id, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *store) Insert(e Entity) error {
	return s.base.Insert(e.GetID(), e)
}

func (s *store) Update(e Entity) error {
	return s.base.Update(e.GetID(), e)
}

func (s *store) Delete(id interface{}) error {
	return s.base.Delete(id)
}

func (s *store) Exists(id interface{}) (bool, error) {
	if se, ok := s.base.(StoreExists); ok {
		return se.Exists(id)
	}

	var _, err = s.Get(id)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, ErrNotFound) {
		return false, nil
	}

	return false, err
}

func (s *store) List() (interface{}, error) {
	if sl, ok := s.base.(StoreList); ok {
		var v = s.config.newList()

		if err := sl.List(v); err != nil {
			return nil, err
		}

		return indirect(v), nil
	}

	return s.Find(nil)
}

func (s *store) Find(filters map[string]interface{}) (interface{}, error) {
	if sf, ok := s.base.(StoreFind); ok {
		var v = s.config.newList()

		if err := sf.Find(filters, v); err != nil {
			return nil, err
		}

		return indirect(v), nil
	}

	return nil, ErrNotImplemented // FIXME wrap error
}

func (s *store) Upsert(e Entity) error {
	if su, ok := s.base.(StoreUpsert); ok {
		return su.Upsert(e.GetID(), e)
	}

	var ok, err = s.Exists(e.GetID())
	if err != nil {
		return err
	}

	if ok {
		return s.Update(e)
	}

	return s.Insert(e)
}

func (s *store) DeleteAll() error {
	if sda, ok := s.base.(StoreDeleteAll); ok {
		return sda.DeleteAll()
	}

	var l, err = s.List()
	if err != nil {
		return err
	}

	var v = reflect.ValueOf(l)
	// FIXME check this is of expected slice type

	for i := 0; i < v.Len(); i++ {
		if err := s.Delete(v.Index(i).Addr().Interface().(Entity).GetID()); err != nil {
			return err
		}
	}

	return nil
}

func indirect(v interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}
