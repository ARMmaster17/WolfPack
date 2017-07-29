package main

import (
	"log"
	"net"
	"strings"
	"time"
)

var packList []node

func alphaMsx(port string, in chan string, out chan string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", strings.Join([]string{":", port}, ""))
	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()
	// ANC is blocking call. testing using concurrent threads.
	go alphaNetChk(ServerConn)
	go alphaMsgChk(in, out)
	for {

	}
}

func alphaMsgChk(in chan string, out chan string) {
	select {
	case dta, ok := <-in:
		if ok {
			// Data is ready to be read.
			if dta == "LIST PACK" {
				rmsg := ""
				for nitm := range packList {
					rmsg = strings.Join([]string{packList[nitm].identifier, "|", packList[nitm].uri, "|", packList[nitm].lastseenat.String(), ","}, "")
				}
				rmsg = rmsg[:len(rmsg)-1]
				out <- rmsg
			}
		} else {
			// Channel got closed for some reason.
		}
	default:
		// Nothing to read here, moving on.
	}
}

func alphaNetChk(conn *net.UDPConn) {
	buf := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buf)
	// Capture payload so we can process it.
	msgdump := string(buf[0:n])
	log.Println("Received ", msgdump, " from ", addr)
	if strings.Contains(msgdump, "CLIENT JOIN ") {
		cname := strings.Replace(msgdump, "CLIENT JOIN ", "", 0)
		packList = append(packList, node{identifier: cname, uri: addr.String(), lastseenat: time.Now()})
	}

	if err != nil {
		log.Println("Error: ", err)
	}
}
