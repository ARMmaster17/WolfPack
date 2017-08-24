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
		time.Sleep(1 * time.Second)
	}
}

func alphaMsgChk(in chan string, out chan string) {
	for {
		select {
		case dta, ok := <-in:
			if ok {
				// Data is ready to be read.
				if dta == "LIST PACK" {
					rmsg := ""
					if len(packList) > 0 {
						for nitm := range packList {
							rmsg = strings.Join([]string{rmsg, packList[nitm].Identifier, "|", packList[nitm].URI, "|", packList[nitm].Lastseenat.String(), ","}, "")
						}
						rmsg = rmsg[:len(rmsg)-1]
					} else {
						rmsg = "EMPTY"
					}
					out <- rmsg
				} else if strings.Contains(dta, "SERVER BANISH") {
					cname := strings.Replace(dta, "SERVER BANISH ", "", 1)
					removeClient(cname)
				}
			} else {
				// Channel got closed for some reason.
			}
		default:
			// Nothing to read here, moving on.
		}
	}
}

func alphaNetChk(conn *net.UDPConn) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		// Capture payload so we can process it.
		msgdump := string(buf[0:n])
		log.Println("Received ", msgdump, " from ", addr)
		if strings.Contains(msgdump, "CLIENT JOIN ") {
			cname := strings.Replace(msgdump, "CLIENT JOIN ", "", 1)
			packList = append(packList, node{Identifier: cname, URI: addr.String(), Lastseenat: time.Now()})
		} else if strings.Contains(msgdump, "CLIENT LEAVE ") {
			cname := strings.Replace(msgdump, "CLIENT LEAVE ", "", 1)
			removeClient(cname)
		}
		if err != nil {
			log.Println("Error: ", err)
		}
	}
}

func removeClient(ident string) {
	idToRemove := -1
	for index, element := range packList {
		if element.Identifier == ident {
			idToRemove = index
			break
		}
	}
	if idToRemove == -1 {
		log.Println(strings.Join([]string{"[MSX][ERROR] Cannot remove client named ", ident}, ""))
	} else {
		packList = removeFromSlice(packList, idToRemove)
	}
}

func removeFromSlice(s []node, index int) []node {
	s[len(s)-1], s[index] = s[index], s[len(s)-1]
	return s[:len(s)-1]
}
