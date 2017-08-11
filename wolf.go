package main

import (
	"net"
	"strings"
	"time"
)

func wolf(port string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", strings.Join([]string{"127.0.0.1:", port}, ""))

	LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")

	Conn, _ := net.DialUDP("udp", LocalAddr, ServerAddr)

	defer Conn.Close()

	_, _ = Conn.Write([]byte("CLIENT JOIN testnode"))
	time.Sleep(4 * time.Second)
	for {
		// Infinitely sleep for...reasons.
		time.Sleep(1 * time.Second)
	}
	//_, _ = Conn.Write([]byte("CLIENT LEAVE testnode"))
}
