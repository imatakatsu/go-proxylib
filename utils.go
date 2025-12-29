package proxylib

import (
	"bytes"
	"net"
	"time"
)

// auth must be "Basic <b64>"
func http_connect(conn net.Conn, target_host string, auth string, deadline time.Time) error {
	buf := make([]byte, 0, 1024)

	buf = append(buf, "CONNECT "...)
	buf = append(buf, target_host...)
	buf = append(buf, " HTTP/1.1\r\nHost: "...)
	buf = append(buf, target_host...)
	buf = append(buf, "\r\n"...)
	if auth != "" {
		buf = append(buf, "Proxy-Authorization: "...)
		buf = append(buf, auth...)
		buf = append(buf, "\r\n"...)
	}
	buf = append(buf, "\r\n"...)

	if !deadline.IsZero() {
		conn.SetDeadline(deadline)
		defer conn.SetDeadline(time.Time{})
	}

	if _, err := conn.Write(buf); err != nil {
		return err
	}

	buf = buf[:cap(buf)]
	tb := 0

	for {
		n, err := conn.Read(buf[tb:])
		if err != nil {
			return err
		}
		tb += n

		if tb >= 4 && bytes.HasSuffix(buf[:tb], []byte("\r\n\r\n")) {
			break
		}

		if tb >= len(buf) {
			return ErrResponseTooLarge
		}
	}

	lineEnd := bytes.IndexByte(buf[:tb], '\r')
	if lineEnd == -1 {
		return ErrInvalidProxyResponse
	}

	if !bytes.Contains(buf[:lineEnd], []byte(" 200 ")) {
		return ErrInvalidProxyResponse
	}

	return nil
}