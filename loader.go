package proxylib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// the returned slice may contain valid proxies
// even if an error is returned
func LoadFromReader(r io.Reader, proto string, parser ParserFunc) ([]*Proxy, error) {
	var proxies []*Proxy

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		p, err := parser(proto, line)
		if err != nil {
			continue
		}

		proxies = append(proxies, p)
	}

	if err := scanner.Err(); err != nil {
		return proxies, err
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies found")
	}
	return proxies, nil
}

func LoadFromFile(path string, proto string, parser ParserFunc) ([]*Proxy, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return LoadFromReader(file, proto, parser)
}

func LoadFromBytes(data []byte, proto string, parser ParserFunc) ([]*Proxy, error) {
	return LoadFromReader(bytes.NewReader(data), proto, parser)
}

func LoadFromString(data string, proto string, parser ParserFunc) ([]*Proxy, error) {
	return LoadFromReader(strings.NewReader(data), proto, parser)
}