package scanner

import (
	"net"
	"strconv"

	"github.com/thbraz/goscan/internal/cli"
)

const (
	dialNetworkTCP = "tcp"
)

type DialResult struct {
	Port int
	Open bool
}

func Scan(args cli.Args) []DialResult {
	var results []DialResult

	for port := range 65536 {
		results = append(results, conn(args.Host, port))
	}

	return results
}

func conn(host string, port int) DialResult {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.Dial(dialNetworkTCP, addr)
	if err != nil {
		return DialResult{Port: port, Open: false}
	}
	conn.Close()

	return DialResult{Port: port, Open: true}
}

