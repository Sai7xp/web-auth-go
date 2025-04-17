package main

import (
	"fmt"
	"log"
	"os"
	internal "web-auth-go/internal"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No option provided to run the code")
	}

	switch args[0] {
	case "1":
		internal.RunJson()
	case "2":
		internal.RunHashing()
	case "3":
		internal.RunHMAC()
	case "4":
		internal.RunJWT()
	default:
		fmt.Println("Invalid Option")
	}
}
