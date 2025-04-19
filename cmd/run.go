package run

import (
	"fmt"
	"log"
	"os"
	internal "web-auth-go/internal"
	passwordAuth "web-auth-go/internal/SessionBasedAndJWTAuthenticationExercise/01-PasswordBasedAuth"
	sessionAuth "web-auth-go/internal/SessionBasedAndJWTAuthenticationExercise/02-SessionBasedAuth"
	JwtTokenBasedAuth "web-auth-go/internal/SessionBasedAndJWTAuthenticationExercise/03-JwtBasedAuth"
)

func Run() {
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
	case "5":
		internal.RunJwtAuthAPI()
	case "6":
		internal.RunBase64()
	case "7":
		passwordAuth.RunPasswordBasedAuth()
	case "8":
		sessionAuth.RunSessionBasedAuth()
	case "9":
		JwtTokenBasedAuth.RunJwtBasedAuth()
	default:
		fmt.Println("Invalid Option")
	}
}
