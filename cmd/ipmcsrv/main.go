package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func tcpClient(c net.Conn) {
	log.Printf("TCP accepted %s\n", c.RemoteAddr().String())
	defer c.Close()

	var out []byte
	in := make([]byte, 4096)

	for {
		if n, err := c.Read(in); err != nil {
			if err == io.EOF {
				log.Printf("TCP %s: EOF", c.RemoteAddr().String())
			} else {
				log.Printf("TCP %s: read error: %v\n", c.RemoteAddr().String(), err)
			}
			break
		} else {
			out = append(out, in[:n]...)
		}
	}

	n, _ := c.Write(out)
	log.Printf("TCP %s: wrote %d bytes, payload is %q\n", c.RemoteAddr().String(), n, string(out[:n]))
}

func tcpServer() {
	addr := cmp.Or(os.Getenv("LISTEN_ADDR"), ":4000")
	l, err := net.Listen("tcp4", addr)
	check(err)
	defer l.Close()

	log.Printf("TCP listening on %q", addr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go tcpClient(c)
	}
}

func udpClient(pc net.PacketConn, addr net.Addr, buf []byte) {
	// Just echo back
	log.Printf("UDP %s: wrote %d bytes, payload is %q\n", addr.String(), len(buf), string(buf))
	pc.WriteTo(buf, addr)
}

func udpUnicastServer() {
	addr := cmp.Or(os.Getenv("UNICAST_LISTEN_PORT"), ":4001")
	pc, err := net.ListenPacket("udp4", addr)
	check(err)

	defer pc.Close()

	log.Printf("UDP listening on %q", addr)

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Printf("UDP %s: %v", addr.String(), err)
			continue
		}
		log.Printf("UDP %s: got %d bytes\n", addr.String(), n)

		go udpClient(pc, addr, buf[:n])
	}
}

func ipmcServer() {
	seqNo := 1

	addrStr := cmp.Or(os.Getenv("MULTICAST_ADDR"), "224.0.0.1:9999")

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	check(err)

	c, err := net.DialUDP("udp", nil, addr)
	check(err)

	defer c.Close()

	log.Printf("IPMC on %s", addr.String())

	for {
		seqNoStr := strconv.Itoa(seqNo)
		c.Write([]byte(seqNoStr))

		time.Sleep(1 * time.Second)

		seqNo += 1
	}

}

func main() {
	ifs, _ := net.Interfaces()
	for idx, i := range ifs {
		fmt.Println("if", idx, i)
		mca, _ := i.MulticastAddrs()
		ifa, _ := i.Addrs()
		for _, a := range ifa {
			fmt.Println("\taddr:", a)
		}

		for _, m := range mca {
			fmt.Println("\tmc:", m)
		}
	}

	go tcpServer()
	go udpUnicastServer()
	ipmcServer()
}
