package main

import (
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/osext"
)

var webInChannel chan string
var webOutChannel chan string

func launchWebGui(webhostCfg string, webportCfg string, in chan string, out chan string) {

	// Set up channel communication.
	webInChannel = in
	webOutChannel = out

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

	// Push out the view with given model information.
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"packlist": getNodeList(),
	})
}

func getNodeList() []node {

	var result []node
	webOutChannel <- "LIST PACK"

	rawResult := <-webInChannel

	nodeArray := strings.Split(rawResult, ",")
	for n := range nodeArray {

		ndta := strings.Split(nodeArray[n], "|")
		lsa, _ := time.Parse(time.RFC3339, ndta[2])
		nd := node{Identifier: ndta[0], URI: ndta[1], Lastseenat: lsa}
		result = append(result, nd)
	}
	return result
}
