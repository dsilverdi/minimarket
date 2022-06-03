package uuid

import (
	"minimarket"
	"minimarket/pkg/errors"

	"github.com/gofrs/uuid"
)

var ErrGeneratingID = errors.New("generating id failed")

var _ minimarket.IDprovider = (*uuidProvider)(nil)

type uuidProvider struct{}

// New instantiates a UUID provider.
func New() minimarket.IDprovider {
	return &uuidProvider{}
}

func (up *uuidProvider) ID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(ErrGeneratingID, err)
	}

	return id.String(), nil
}
