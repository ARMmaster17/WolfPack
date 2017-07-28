package main

import (
	"log"
	"net"
	"strings"
)

func alphaMsx(port string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", strings.Join([]string{":", port}, ""))
	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			log.Println("Error: ", err)
		}
	}
}
