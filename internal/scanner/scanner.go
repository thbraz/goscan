package scanner

import (
	"net"
	"strconv"

)

const (
	dialNetworkTCP = "tcp"
	maxPort = 65536
)

type DialResult struct {
	Port int
	Open bool
}

func Scan(host string) []DialResult {
	results := make([]DialResult, 0, maxPort)

	for port := 1; port <= maxPort; port++ {
		results = append(results, probePort(host, port))
	}

	return results
}

func probePort(host string, port int) DialResult {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.Dial(dialNetworkTCP, addr)
	if err != nil {
		return DialResult{Port: port, Open: false}
	}
	conn.Close()

	return DialResult{Port: port, Open: true}
}