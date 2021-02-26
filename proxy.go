package main

import (
	"fmt"
	"sync"
	"time"
	"os"
	"crypto/tls"
	"net/url"
	"net"
	"net/http"
)

func main() {
	secure := os.Getenv("SECURE")
	pu := os.Getenv("PROXY")
	if len(pu) == 0 {
		pu = "https://localhost:443"
		secure = "1"
	}
	u := os.Getenv("URL")
	if len(u) == 0 {
		fmt.Println("Please set URL properly")
		os.Exit(1)
	}
	prox, err := url.Parse(pu)
	if err != nil {
		fmt.Println("Please setup valid proxy url")
		os.Exit(1)
	}
	var transport *http.Transport = &http.Transport{
		Proxy: http.ProxyURL(prox),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if len(secure) > 0 {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	var wg sync.WaitGroup
	s := time.Now()
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			req, _ := http.NewRequest("GET", u, nil)
			resp, err := transport.RoundTrip(req)
			if err != nil {
				fmt.Printf("Request failed %+v\n", err)
				return
			}
			resp.Body.Close()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Total time: ", time.Since(s))
}
