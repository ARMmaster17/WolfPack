package main

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"fmt"

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
	log.Println("Initializing message exchange")
	outbound := make(chan string)
	inbound := make(chan string)
	log.Println("Launching slave client")
	go wolf(outbound, inbound)
	// Verify client is up
	log.Println("Testing connection")
	outbound <- "ping"
	inbuffer := <-inbound
	if inbuffer == "pong" {
		log.Println("Slave unit connection sucessful")
	}
	//////////////////////////////
	portCfg := os.Getenv("port")
	if portCfg == "" {
		log.Fatalln("$PORT must be set.")
	}
	hostCfg := os.Getenv("host")
	if hostCfg == "" {
		log.Fatalln("$HOST must be set.")
	}
	fmt.Printf("Web GUI available at %v:%v\n", hostCfg, portCfg)
	// Need to use release mode because of issue #119 on gin-gonic/gin.
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Grab our current working directory.
	wd, _ := osext.ExecutableFolder()
	// Preload all templates into memory.
	router.LoadHTMLGlob(path.Join(wd, "templates", "*"))
	// Assign base route to home controller.
	router.GET("/", controllerhome)
	// Everything is set up, let gin do its thing.
	err := router.Run(strings.Join([]string{hostCfg, ":", portCfg}, ""))
	log.Fatalf("Gin error: %v\n", err)
}

func wolf(in chan string, out chan string) {
	for {
		msg := <-in
		if msg == "ping" {
			out <- "pong"
		}
	}
}

type node struct {
	identifier string
	uri        string
	lastseenat time.Time
}
