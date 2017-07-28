package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
)

func launchWebGui(webhostCfg string, webportCfg string) {
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
