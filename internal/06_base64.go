package internal

import (
	"encoding/base64"
	"fmt"
)

func RunBase64() {
	// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
	// var StdEncoding = NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	// URLEncoding is the alternate base64 encoding defined in RFC 4648.
	// It is typically used in URLs and file names.
	// var URLEncoding = NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

	base64URL := base64.URLEncoding.EncodeToString([]byte("Hello ./ and %$😎"))
	fmt.Println(base64URL)

	decodedMessage, _ := base64.URLEncoding.DecodeString(base64URL)
	fmt.Println(string(decodedMessage))
}
