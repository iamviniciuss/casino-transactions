package http

import (
	"context"
	"net"
)

type FiberHandlerFunc func(context.Context, map[string]string, []byte, QueryParams, LocalsFunc) (interface{}, *IntegrationError)

type LocalsFunc func(key interface{}) interface{}

type HttpService interface {
	Get(path string, callback FiberHandlerFunc)
	Post(path string, callback FiberHandlerFunc)
	ListenAndServe(port string) error
	Listen(listener net.Listener) error
}

type QueryParams interface {
	GetParam(key string) []byte
	AddParam(key string, value string)
}
