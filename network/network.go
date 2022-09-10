package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println(GetOutboundIP())

	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	go func() {
		addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8829")
		if err != nil {
			panic(err)
		}

		_, err = pc.WriteTo([]byte("data to transmit"), addr)
		if err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s sent this: %s\n", addr, buf[:n])
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
