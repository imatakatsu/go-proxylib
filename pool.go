package proxylib

import (
	"github.com/imatakatsu/ringer"
)

type Pool struct {
	*ringer.Rotator[*Proxy]
}

func NewPool(proxies []*Proxy) (*Pool, error) {
	if len(proxies) < 1 {
		return nil, ErrEmptyProxyList
	}
	return &Pool{
		ringer.NewRotator(proxies),
	}, nil
}