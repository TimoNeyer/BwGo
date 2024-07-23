package main

import "fmt"

func main() {
	BwServer := Server{
		hostname: "localhost",
		port:     8087,
		unlocked: -1,
		Token:    "",
	}
	err := BwServer.unlock_server()
	if err != nil {
		fmt.Printf("Error:\n%s", err.Error())
	}

}
