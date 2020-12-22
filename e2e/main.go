package main

import "time"

func main() {
	socket()

	time.Sleep(1 * time.Second)

	// Need
	// - mqtt broker server
	// - mosquitto_pub
	broker()
}
