package minimarket

type IDprovider interface {
	ID() (string, error)
}
