package test

import (
	"net/url"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/mock"
)

type MockedCall[I, R any] struct {
	Input  I
	Return R
	Error  error
}

type MockSite struct {
	mock.Mock
}

func (s *MockSite) GetInput(url *url.URL, v ...string) (*model.SiteInput, error) {
	args := s.Called(url, v)
	return args.Get(0).(*model.SiteInput), args.Error(1)
}

func (s *MockSite) Fetch(fi model.SiteInput) (model.Posts, error) {
	args := s.Called(fi)
	return args.Get(0).(model.Posts), args.Error(1)
}
