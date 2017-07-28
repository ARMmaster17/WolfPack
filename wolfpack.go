package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
)

func controllerhome(c *gin.Context) {
	// Grab the machine hostname.
	hn, _ := os.Hostname()
	// Test the node struct by utilizing it in view/model bindings.
	testnode := node{identifier: hn, uri: "http://localhost:8080/", lastseenat: time.Now()}
	// Push out the view with given model information.
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"hostname":    testnode.identifier,
		"lastupdated": time.Now().String(),
		"wolfcount":   5,
	})
}

func main() {
	//////////////////////////////
	// DEBUG: in the future this will
	// happen on a separate machine over UDP
	//log.Println("Initializing message exchange")
	//outbound := make(chan string)
	//inbound := make(chan string)
	//log.Println("Launching slave client")
	//go wolf(outbound, inbound)
	// Verify client is up
	//log.Println("Testing connection")
	//outbound <- "ping"
	//inbuffer := <-inbound
	//if inbuffer == "pong" {
	//	log.Println("Slave unit connection sucessful")
	//}
	//////////////////////////////
	webhostCfg := setConfigValue("WOLF_WEBHOST", "0.0.0.0")
	webportCfg := setConfigValue("WOLF_WEBPORT", "8080")
	msxportCfg := setConfigValue("WOLF_MSXPORT", "9300")
	// Launch MSX controller.
	go alphaMsx(msxportCfg)
	time.Sleep(1 * time.Second)
	go wolf(msxportCfg)
	// Need to use release mode because of issue #119 on gin-gonic/gin.
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Grab our current working directory.
	wd, _ := osext.ExecutableFolder()
	// Preload all templates into memory.
	router.LoadHTMLGlob(path.Join(wd, "templates", "*"))
	// Assign base route to home controller.
	router.GET("/", controllerhome)
	// Notify user of connection credentials to web GUI.
	log.Printf("Web GUI available at %v:%v\n", webhostCfg, webportCfg)
	// Everything is set up, let gin do its thing.
	err := router.Run(strings.Join([]string{webhostCfg, ":", webportCfg}, ""))
	log.Fatalf("Gin error: %v\n", err)
}

func setConfigValue(parameter string, defaultValue string) string {
	value := os.Getenv(parameter)
	if value == "" {
		log.Println(strings.Join([]string{"$", parameter, " not set. Using program defaults."}, ""))
		return defaultValue
	}
	return value
}

func alphaMsx(port string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", strings.Join([]string{":", port}, ""))
	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func wolf(port string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", strings.Join([]string{"127.0.0.1:", port}, ""))

	LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")

	Conn, _ := net.DialUDP("udp", LocalAddr, ServerAddr)

	defer Conn.Close()

	_, _ = Conn.Write([]byte("CLIENT JOIN testnode"))
	time.Sleep(4 * time.Second)
	_, _ = Conn.Write([]byte("CLIENT LEAVE testnode"))
}

type node struct {
	identifier string
	uri        string
	lastseenat time.Time
}
