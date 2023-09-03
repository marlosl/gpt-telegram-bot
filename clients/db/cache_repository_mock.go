package db

import (
	"github.com/stretchr/testify/mock"
)

type CacheRepositorMock struct {
	mock.Mock
}

func (c *CacheRepositorMock) SaveItem(item *string) error {
	args := c.Called(item)
	return args.Error(0)
}

func (c *CacheRepositorMock) ItemExists(item string) bool {
	args := c.Called(item)
	return args.Bool(0)
}

func (c *CacheRepositorMock) GetItem(item string) (*string, error) {
	args := c.Called(item)
	return args.Get(0).(*string), args.Error(1)
}

func (c *CacheRepositorMock) DeleteItem(item string) error {
	args := c.Called(item)
	return args.Error(0)
}
