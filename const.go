package proxylib

import (
	"errors"
)

const (
	HTTP = "http"
	HTTPS = "https"
	SOCKS5 = "socks5"
	SOCKS4 = "socks4"
)

var (
	ErrResponseTooLarge = errors.New("response too large")
	ErrInvalidProxyResponse = errors.New("got invalid proxy response")
	ErrNoProxiesFound = errors.New("no proxies found")
	ErrInvalidProxyFormat = errors.New("invalid proxy format provided")
	ErrNotSupported = errors.New("this protocol is not supported yet, use HTTP/HTTPS instead :(")
	ErrEmptyProxyList = errors.New("provided proxy list is emptu")
)