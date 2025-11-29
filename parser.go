package proxylib

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type ParserFunc func(proto string, line string) (*Proxy, error)

// valid url format
// http://[login:[password]@]host:port
// proto can be null
func ParseURL(proto string, line string) (*Proxy, error) {
	var (
		username string
		password string
		host     string
		port     string
	)

	u, err := url.Parse(line)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "" {
		proto = u.Scheme
	}

	if proto == "" {
		return nil, fmt.Errorf("proxy proto not found")
	}

	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	host = u.Hostname()
	if u.Port() == "" {
		return nil, fmt.Errorf("invalid url format, port MUST exist")
	}
	port = u.Port()
	return &Proxy{
		Protocol: proto,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		b64Auth:  "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
	}, nil
}

// pseudo url format
// [login:[password]@]host:port
// proto MUST be not null
func ParsePseudoURL(proto string, line string) (*Proxy, error) {
	return ParseURL("", proto+"://"+line)
}

// string proxy format
// host:port[:login[:password]]
// proto MUST be not null
func ParseString(proto string, line string) (*Proxy, error) {
	var (
		username string
		password string
		host     string
		port     string
	)
	parts := strings.SplitN(line, ":", 4)

	switch len(parts) {
	case 2:
		host = parts[0]
		port = parts[1]
	case 3:
		host = parts[0]
		port = parts[1]
		username = parts[2]
	case 4:
		host = parts[0]
		port = parts[1]
		username = parts[2]
		password = parts[3]
	default:
		return nil, fmt.Errorf("invalid proxy string format")
	}

	return &Proxy{
		Protocol: proto,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		b64Auth:  "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
	}, nil
}
