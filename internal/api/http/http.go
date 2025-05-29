package http

import (
	"context"
	"net"
)

type LocalsFunc func(key interface{}) interface{}

type HttpService interface {
	Get(path string, callback func(context.Context, map[string]string, []byte, QueryParams, LocalsFunc) (interface{}, *IntegrationError))
	Post(path string, callback func(context.Context, map[string]string, []byte, QueryParams, LocalsFunc) (interface{}, *IntegrationError))
	ListenAndServe(port string) error
	Listen(listener net.Listener) error
}

type QueryParams interface {
	GetParam(key string) []byte
	AddParam(key string, value string)
}
