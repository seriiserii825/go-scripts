package main

import "github.com/seriiserii825/go-libs/msg"

func main() {
	errorMessage := "An error occurred"
	msg.Success("Hello, World!")
	msg.Error(errorMessage)
}
