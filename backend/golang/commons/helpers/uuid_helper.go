package helpers

import "github.com/google/uuid"

type UuidHelper interface {
	String() string
}

type UuidHelperImplementation struct {
}

func NewUuidHelper() UuidHelper {
	return &UuidHelperImplementation{}
}

func (helperImplementation *UuidHelperImplementation) String() string {
	return uuid.New().String()
}
