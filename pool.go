package proxylib

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	proxies []*Proxy
	mu sync.RWMutex
	idx     uint64
}

func (p *Pool) Next() *Proxy {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.proxies) < 1 {
		return &Proxy{}
	}

	idx := atomic.AddUint64(&p.idx, 1)

	return p.proxies[int(idx) % len(p.proxies)]
}

/*
	functions below will be added in further updates
	you can donate me if u want to see them in proxylib faster ;d
	
*/

/*
type ProxyItem struct {
	*Proxy

	freezeUntil int64
	inUse       int32
	isBorrowed  int32
	failCounter int32
	isRemoved   int32
}

type Conn struct {
	net.Conn
	inUse *int32
	Borrowed *int32
	once sync.Once
}

func (c *Conn) Close() error {
	c.once.Do(func() {
		atomic.StoreInt32(c.Borrowed, 0)
		atomic.AddInt32(c.inUse, -1)
	})

	return c.Conn.Close()
}

func (p *ProxyItem) FreezeFor(duration time.Duration) {
	atomic.StoreInt64(&p.freezeUntil, time.Now().Add(duration).UnixMilli())
}

func (p *ProxyItem) FreezeUntil(date time.Time) {
	atomic.StoreInt64(&p.freezeUntil, date.UnixMilli())
}

func (p *ProxyItem) DialTimeout(host string, timeout time.Duration) (*Conn, error) {
	atomic.StoreInt32(&p.isBorrowed, 1)
	atomic.AddInt32(&p.inUse, 1)

	net_conn, err := p.Proxy.DialTimeout(host, timeout)
	if err != nil {
		atomic.AddInt32(&p.failCounter, 1)
		atomic.AddInt32(&p.inUse, -1)
		atomic.StoreInt32(&p.isBorrowed, 0)
		return nil, err
	}

	atomic.StoreInt32(&p.failCounter, 0)

	return &Conn{
		net_conn,
		&p.inUse,
		&p.isBorrowed,
		sync.Once{},
	}, nil
}

func (p *ProxyItem) Dial(host string) (*Conn, error) {
	return p.DialTimeout(host, 0)
}
*/
