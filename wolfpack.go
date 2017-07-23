package main

import (
	"net/http"
	"os"
	"path"
	"time"

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
	})
}

func main() {
	// Need to use release mode because of issue #119 on gin-gonic/gin.
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Grab our current working directory.
	wd, _ := osext.ExecutableFolder()
	// This is needed for debugging because VSCode likes to be difficult.
	//wd := path.Join("C:", "Users", "wordm", "go", "src", "WolfPack")
	// Preload all templates into memory.
	router.LoadHTMLGlob(path.Join(wd, "templates", "*"))
	// Assign base route to home controller.
	router.GET("/", controllerhome)
	// Everything is set up, let gin do its thing.
	router.Run(":8080")
}

type node struct {
	identifier string
	uri        string
	lastseenat time.Time
}
