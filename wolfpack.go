package main

import "time"

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
	envVars := loadEnvironmentVars()
	// Launch MSX controller.
	go alphaMsx(envVars["msx_port"])
	time.Sleep(1 * time.Second)
	go wolf(envVars["msx_port"])
	launchWebGui(envVars["web_host"], envVars["web_port"])
}
