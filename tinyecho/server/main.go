package main

import (
	"flag"
	"fmt"
	"tinyecho/server/tcpserver"
	"tinyecho/server/udpserver"
)

var (
	protocol string
	address  string
)

func init() {
	flag.StringVar(&protocol, "proto", "tcp", "network protocol")
	flag.StringVar(&address, "addr", "127.0.0.1:8897", "address")
	flag.Parse()
}

func main() {
	switch protocol {
	case tcpserver.Protocol:
		tcpserver.Run(address)
	case udpserver.Protocol:
		udpserver.Run(address)
	default:
		fmt.Println("Server protocol error: ", protocol)
	}
}
