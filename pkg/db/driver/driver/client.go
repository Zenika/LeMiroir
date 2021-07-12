package driver

type Client interface {
	Store(...StoreOption) (Store, error)
	Begin() (Transaction, error)
	Close() error
}
