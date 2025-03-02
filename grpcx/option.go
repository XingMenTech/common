package grpcx

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/resolver"
	"time"
)

// ServerOption service side option
type ServerOption func(*serverOptions)

type serverOptions struct {
	httpAddr  string
	httpSetup func(*gin.Engine)
}

// ClientOption is client create option
type ClientOption func(*clientOptions)

type clientOptions struct {
	prefix   string
	resolver resolver.Builder
	timeout  time.Duration
}

// WithDirectAddresses returns a direct addresses option
func WithDirectAddresses(addrs ...string) ClientOption {
	return func(opts *clientOptions) {
		if len(addrs) <= 0 {
			return
		}
		opts.prefix = addrs[0]
		//opts.resolver = resolver.Builder(addrs...)
	}
}

// WithTimeout returns a timeout option
func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *clientOptions) {
		opts.timeout = timeout
	}
}
