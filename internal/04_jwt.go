package internal

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func RunJWT() {
	newToken := generateToken()

	time.Sleep(time.Second * 8)

	claims, err := verifyToken(newToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Valid Token. Data: ", claims)
}

func generateToken() string {
	claims := jwt.MapClaims{
		"username": "sai7xp",
		"exp":      time.Now().Add(time.Second * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatal("Error while generating the token: ", err)
	}
	fmt.Println("JWT: ", tokenString)
	return tokenString
}

func verifyToken(token string) (jwt.MapClaims, error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return secretKey, nil })
	if err != nil {
		log.Fatal(err)
	}
	if !tkn.Valid {
		return nil, errors.New("token expired")
	}
	return tkn.Claims.(jwt.MapClaims), nil
}
