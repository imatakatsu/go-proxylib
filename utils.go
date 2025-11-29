package proxylib

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

// auth must be "Basic <b64>"
func http_connect(conn net.Conn, target_host string, auth string, deadline time.Time) error {
	if auth != "" {
		auth = "Proxy-Authorization: " + auth + "\r\n"
	}
	_, err := fmt.Fprintf(conn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n%s\r\n", target_host, target_host, auth)
	if err != nil {
		return err
	}

	var buf [4096]byte
	if !deadline.IsZero() {
		conn.SetDeadline(deadline)
		defer conn.SetDeadline(time.Time{})
	}

	totbytes := 0

	for {
		n, err := conn.Read(buf[totbytes:])
		if err != nil {
			return err
		}

		totbytes += n

		if bytes.HasSuffix(buf[:totbytes], []byte("\r\n\r\n")) {
			break
		} else if totbytes == len(buf) {
			return fmt.Errorf("response too large")
		}
	}

	parts := bytes.SplitN(buf[:totbytes], []byte("\r\n"), 2)
	if !bytes.Contains(parts[0], []byte(" 200 ")) {
		return fmt.Errorf("bad response status from server: %s", parts[0])
	}

	return nil
}
