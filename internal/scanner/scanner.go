package scanner

import (
	"net"
	"strconv"
	"sync"
)

const (
	dialNetworkTCP = "tcp"
	maxPort = 65536
	minPort = 1
)

type DialResult struct {
	Port int
	Open bool
}

func Scan(host string) []DialResult {
	results := make([]DialResult, maxPort)
	sem := make(chan struct{}, 500)

	var wg sync.WaitGroup
	for port := minPort; port <= maxPort; port++ {
		sem <- struct{}{}
		
		wg.Go(func() {
			defer func() { <-sem }()
			results[port-1] = probePort(host, port)
		})
	}
	wg.Wait()

	return results
}

func probePort(host string, port int) DialResult {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.Dial(dialNetworkTCP, addr)
	open := err == nil
	if open {
		conn.Close()
	}

	return DialResult{Port: port, Open: open}
}