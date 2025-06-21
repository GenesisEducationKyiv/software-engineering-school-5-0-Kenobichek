package httpclient

import (
	"net"
	"net/http"
	"time"
)

const (
	defaultRequestTimeout        = 5 * time.Second
	defaultMaxIdleConns          = 100
	defaultMaxIdleConnsPerHost   = 10
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultDialTimeout           = 5 * time.Second
	defaultKeepAlive             = 30 * time.Second
)

func New() *http.Client {
	return &http.Client{
		Timeout:   defaultRequestTimeout,
		Transport: newTransport(),
	}
}

func newTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          defaultMaxIdleConns,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout,
			KeepAlive: defaultKeepAlive,
		}).DialContext,
	}
}
