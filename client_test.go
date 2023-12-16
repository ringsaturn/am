package am_test

import (
	"context"
	"sync"
)

type FooTokenSaver struct {
	token string
	exp   int64
	mutex sync.Mutex
}

func (s *FooTokenSaver) GetAccessToken(context.Context) (string, int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.token, s.exp, nil
}

func (s *FooTokenSaver) SetAccessToken(ctx context.Context, token string, exp int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.token = token
	s.exp = exp
	return nil
}
