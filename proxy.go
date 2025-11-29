package proxylib

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type Proxy struct {
	Protocol string
	Host     string
	Port     string
	Username string
	Password string
	b64Auth  string
}

func (p *Proxy) Dial(host string) (net.Conn, error) {
	return p.DialTimeout(host, 0)
}

func (p *Proxy) DialTimeout(host string, timeout time.Duration) (net.Conn, error) {
	var deadline time.Time

	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(p.Host, p.Port), timeout)
	if err != nil {
		return nil, err
	}

	switch p.Protocol {
	case HTTP:
		err := http_connect(conn, host, p.b64Auth, deadline)
		if err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil

	case HTTPS:
		tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true, ServerName: p.Host})

		if !deadline.IsZero() {
			conn.SetDeadline(deadline)
		}

		err := tlsConn.Handshake()
		if err != nil {
			tlsConn.Close()
			return nil, err
		}

		err = http_connect(tlsConn, host, p.b64Auth, deadline)
		if err != nil {
			tlsConn.Close()
			return nil, err
		}
		return tlsConn, nil

	default:
		conn.Close()
		return nil, fmt.Errorf("this protocol is not supported yet, use HTTP/HTTPS instead :(")
	}
}
