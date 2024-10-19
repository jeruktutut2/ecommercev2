package mockhelpers

import "github.com/stretchr/testify/mock"

type UuidHelperMock struct {
	Mock mock.Mock
}

func (helper *UuidHelperMock) String() string {
	arguments := helper.Mock.Called()
	return arguments.Get(0).(string)
}
