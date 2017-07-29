package main

import "time"
import "log"

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

	// Load in environment variables.
	log.Println("Loading environment variables...")
	envVars := loadEnvironmentVars()
	// Set up communication channels between MSX and WGUI processes.
	log.Println("Initializing inter-process pipes...")
	MtW := make(chan string)
	WtM := make(chan string)
	// Launch MSX controller.
	log.Println("Launching MSX controller...")
	go alphaMsx(envVars["msx_port"], WtM, MtW)
	// Wait one second while that gets set up.
	time.Sleep(1 * time.Second)
	// Launch a bunch of wolf nodes as golang processes.
	log.Println("Launching slave processes...")
	go wolf(envVars["msx_port"])
	go wolf(envVars["msx_port"])
	go wolf(envVars["msx_port"])
	// Launch WebGUI. This will be the blocking process until WolfPack terminates.
	log.Println("Launching WebGUI...")
	launchWebGui(envVars["web_host"], envVars["web_port"], MtW, WtM)
}
