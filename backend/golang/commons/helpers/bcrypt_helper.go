package helpers

import "golang.org/x/crypto/bcrypt"

type BcryptHelper interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type BycryptHelperImplementation struct {
}

func NewBcryptHelper() BcryptHelper {
	return &BycryptHelperImplementation{}
}

func (helper *BycryptHelperImplementation) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
