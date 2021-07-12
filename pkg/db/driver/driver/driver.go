package driver

import (
	"fmt"
)

var drivers = make(map[string]Driver)

type Driver interface {
	Open() (Client, error)
}

func Register(name string, d Driver) {
	drivers[name] = d
}

func Get(name string) Driver {
	if d, ok := drivers[name]; ok {
		return d
	}
	panic(fmt.Errorf("Unknown database driver %s", name))
}
