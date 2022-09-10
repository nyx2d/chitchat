package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const advertisePort = 8009

var packetConn net.PacketConn

type peerInfo struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	// TODO: include public key too
}

func initNetwork() {
	var err error
	packetConn, err = net.ListenPacket("udp4", fmt.Sprintf(":%d", advertisePort))
	if err != nil {
		panic(err)
	}
}

// Advertise ...
func advertise() {
	nodeAddress, err := getNodeAddress()
	if err != nil {
		panic(err)
	}

	// TODO: use subnet broadcast address, support ipv6, etc
	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", advertisePort))
	if err != nil {
		panic(err)
	}

	info := peerInfo{
		Address: nodeAddress.String(),
		Port:    "8010", // TODO: switch with actual conn port
		// TODO: populate public key
	}
	serializedInfo, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Second) {
		_, err := packetConn.WriteTo(serializedInfo, broadcastAddr)
		if err != nil {
			panic(err)
		}
	}
}

func findPeers() {
	for {
		buf := make([]byte, 2048)
		n, addr, err := packetConn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		var info peerInfo
		err = json.Unmarshal(buf[:n], &info)
		if err != nil {
			log.Println(fmt.Sprintf("bad info payload recieved from %s", addr))
			continue
		}
		fmt.Println(info)
	}
}

func getNodeAddress() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
