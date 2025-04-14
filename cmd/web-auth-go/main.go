package main

import (
	"fmt"
	"log"
	"os"
	myjson "web-auth-go/internal/01-json"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No option provided to run the code")
	}

	switch args[0] {
	case "1":
		myjson.RunJson()
	default:
		fmt.Println("Invalid Option")
	}
}
