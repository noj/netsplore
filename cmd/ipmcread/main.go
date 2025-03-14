package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//ifArg := flag.String("interface", "", "interface to bind to")
	mcAddrArg := flag.String("addr", "224.0.0.1:9999", "interface to bind to")
	flag.Parse()

	ifs, _ := net.Interfaces()
	for idx, i := range ifs {
		fmt.Println("if", idx, i)
		mca, _ := i.MulticastAddrs()
		for _, m := range mca {
			fmt.Println("\t", m)
		}
	}

	ifa, _ := net.InterfaceAddrs()
	for idx, a := range ifa {
		fmt.Println("addr", idx, a)
	}

	// var iface *net.Interface
	// if ifArg != nil {
	// 	//iface = net.Interface
	// 	iface = nil
	// }

	addr, err := net.ResolveUDPAddr("udp", *mcAddrArg)
	check(err)

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	packet := make([]byte, 1500)
	l.SetReadBuffer(len(packet))

	for {
		n, src, err := l.ReadFromUDP(packet)
		check(err)

		log.Printf("%s: %s", src.String(), hex.Dump(packet[:n]))
	}
}
