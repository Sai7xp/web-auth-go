package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func RunHashing() {
	fmt.Println("Hashing")
	bytes := sha256.Sum256([]byte("password"))
	// hashing algorithms do the math on bits and returns the output in raw bytes
	// to convert the raw bytes in printable format, we have to convert those raw bytes to either
	// - hexadecimal
	// - base64

	hexString := hex.EncodeToString(bytes[:]) // bs[:] trick to convert array to slice
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
