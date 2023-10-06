package am_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ringsaturn/am"
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

func ExampleNewClient() {
	simpleClient := am.NewClient("your_auth_token")
	fmt.Println(simpleClient)

	customSaverClient := am.NewClient("your_auth_token", am.WithTokenSaver(&FooTokenSaver{}))
	fmt.Println(customSaverClient)

	customHTTP := &http.Client{
		Timeout: 10,
	}
	customHTTPClient := am.NewClient("your_auth_token", am.WithHTTPClient(customHTTP))
	fmt.Println(customHTTPClient)

}
