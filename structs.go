package ipGeo

import "net"

type IPv4DB struct {
	data         []byte
	IndexStart   uint32
	RecordsCount uint32
	ipSize       uint32
	refSize      uint32
	indexSize    uint32
}

type IPv6DB struct {
	data         []byte
	IndexStart   uint64
	RecordsCount uint64
	ipSize       uint64
	refSize      uint64
	indexSize    uint64
}

type Result struct {
	IP      net.IP
	Country string
	Area    string
}

type IPDB interface {
	GetIPLocation(string) (*Result, error)
	GetIPNumLocation(net.IP) (*Result, error)
}
