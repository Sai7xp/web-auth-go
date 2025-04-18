package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func RunHashing() {
	fmt.Println("Hashing")
	hexBytes := sha256.Sum256([]byte("password"))
	hexString := hex.EncodeToString(hexBytes[:]) // bs[:] trick to convert array to slice
	fmt.Println(hexString)
	/*
		You can do the same thing in terminal as well

		// hash of a string
		`echo -n "password" | shasum -a 256`

		// generate the hash of a file
		shasum -a 256 go1.24.2.darwin-amd64.pkg
		// you can verify the sha256 checksum on go lang downloads page
	*/

	// hashing passwords along with the salt
	// bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	// fmt.Println(string(bcryptBytes))
}
