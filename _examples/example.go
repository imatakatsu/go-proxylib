package main

import (
	"fmt"
	"log"
	"time"

	"github.com/imatakatsu/go-proxylib"
)

/*
this is proxies.txt file
content example with
invalid proxies

187.214.66.321:999
11.55.432.42:8080
102.210.105.351:83
131.257.91.113:999
45.22.299.127:8888
80.913.106.129:8081
195.142.258.190:8080
103.54.811.152:8085
*/

func main() {
	proxies, err := proxylib.LoadFromFile("proxies.txt", proxylib.HTTP, proxylib.ParseString)
	// if parsed proxy list length == 0
	// nil slice will be returned
	if err != nil && proxies == nil {
		panic("failed to parse proxies")
	}
	fmt.Println("proxies parsed!")


	/* 
		default proxy usage example
								     */
	for _, p := range proxies {
		log.Printf("connecting to %s:%s...\r\n", p.Host, p.Port)
		conn, err := p.DialTimeout("ident.me:80", time.Second*time.Duration(10))
		if err != nil {
			log.Printf("invalid proxy got, err: %s\r\n", err.Error())
			continue
		}
		defer conn.Close()

		_, err = conn.Write([]byte("GET ident.me:80 HTTP/1.1\r\nHost: ident.me:80\r\nConnection: close\r\n\r\n"))
		if err != nil {
			log.Printf("failed to send request via proxy, err: %s\r\n", err.Error())
			continue
		}

		var buf [4096]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("failed to read answer from server, err: %s\r\n", err.Error())
			continue
		}
		fmt.Println(string(buf[:n]))
		break // break on first success
	}
	fmt.Println("all proxies checked")


	/*
		proxy pool usage example
								  */
	var pool proxylib.Pool
	pool.AddProxies(proxies)
	for {
		// return next proxy using round robin
		// Next() is safe for concurrent threads
		p := pool.Next()
		conn, err := p.DialTimeout("ident.me:80", time.Second*time.Duration(10))
		if err != nil {
			log.Printf("invalid proxy got, err: %s\r\n", err.Error())
			continue
		}
		defer conn.Close()

		_, err = conn.Write([]byte("GET ident.me:80 HTTP/1.1\r\nHost: ident.me:80\r\nConnection: close\r\n\r\n"))
		if err != nil {
			log.Printf("failed to send request via proxy, err: %s\r\n", err.Error())
			continue
		}

		var buf [4096]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("failed to read answer from server, err: %s\r\n", err.Error())
			continue
		}
		fmt.Println(string(buf[:n]))
		break // break on first success
	}
	fmt.Println("valid proxy found!!")
}
