package main

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func get_password() (password []byte, err error) {
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte{0}, fmt.Errorf("Failed to get Password")
	}
	return bytepw, nil
}
