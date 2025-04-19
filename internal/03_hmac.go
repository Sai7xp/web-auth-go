package internal

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// Hash-based Message Authentication Code
// Sender will have a message to transfer and a shared secret key to generate the hmac
// hmac = Hash(message + <secret_key>) this will be sent to receiver along with the message.
// Receiver will generate the hmac code again with the secret key and compares the new signature and the one received from sender
func RunHMAC() {

	message := "my-message"
	secretKey := "symm-key"
	signature := generateHmac(message, secretKey)

	// Send message
	SendToReceiver(message, signature)
}

// Sender side steps:
// returns a signature of given message
func generateHmac(msg, symmKey string) string {
	// pass the hashing algorithm that you want to use and secret key. hmac generates the hash combining both
	h := hmac.New(sha512.New, []byte(symmKey))
	h.Write([]byte(msg))
	signatureBytes := h.Sum(nil)

	signatureAsHexString := hex.EncodeToString(signatureBytes)
	fmt.Println(signatureAsHexString)
	return signatureAsHexString
}

func SendToReceiver(msg, hmacFromSender string) {
	// we received a message, now verify the signature by generating the hmac again with the same secret key and message
	secretKey := "symm-key"

	hamcGeneratedAtReceiverEnd := generateHmac(msg, secretKey)

	areEqual := hmac.Equal([]byte(hamcGeneratedAtReceiverEnd), []byte(hmacFromSender))
	fmt.Println(areEqual)
}
