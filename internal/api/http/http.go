package http

import (
	"context"
	"net"
)

type FibberHandlerFunc func(context.Context, map[string]string, []byte, QueryParams, LocalsFunc) (interface{}, *IntegrationError) 

type LocalsFunc func(key interface{}) interface{}

type HttpService interface {
	Get(path string, callback FibberHandlerFunc)
	Post(path string, callback FibberHandlerFunc)
	ListenAndServe(port string) error
	Listen(listener net.Listener) error
}

type QueryParams interface {
	GetParam(key string) []byte
	AddParam(key string, value string)
}
