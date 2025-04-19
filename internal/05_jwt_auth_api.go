package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	PORT                     = ":7658"
	userClaimsKey contextKey = "userClaims"
)

var jwt_secret = []byte("secret-key")

func RunJwtAuthAPI() {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", RecoveryHandler(loginHandler))
	mux.HandleFunc("GET /posts", ValidateTokenMiddleware(getUserPostsHandler)) // we can use RecoveryHandler for this as well

	log.Printf("listening on port %s\n", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}

}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Panic("Error while reading r.Body: ", err)
	}

	if creds.Username != "testuser" || creds.Password != "test-jwt-123" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// create token for the user and send it in response
	token := generateJWT(creds.Username)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login Successful", "token": token})

	// or we can simply set the token in cookie
	// and from the next time token will be added to every request that comes from client
	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: token})
}

func getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the user claims that we set during validate token middleware
	claims, ok := r.Context().Value(userClaimsKey).(*UserClaims)
	if !ok {
		http.Error(w, "your claims not found in the token", http.StatusUnauthorized)
		return
	}
	// validate token middleware will be executed before this handler
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello ` + claims.Username + ` Welcome to the protected route. Here is your data" }`))
	w.WriteHeader(http.StatusOK)
}

/*
Util method to generate the JWT token for a user
*/
type UserClaims struct {
	Username             string
	jwt.RegisteredClaims // embedding
}

func generateJWT(username string) string {
	userClaims := UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 10)}, // this expiration minutes will come from env config in production
		},
	}
	// [userClaims] is of type custom struct that we have created - UserClaims
	// we were able to pass our struct to NewWithClaims because we have embedded jwt.RegisteredClaims in our struct.
	// [RegisteredClaims] implements all the methods defined in [jwt.Claims] interface
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	signedString, err := token.SignedString(jwt_secret)
	if err != nil {
		log.Panic("Error while generating the token: ", err)
	}
	return signedString
}

/*
 Common Middleware for token verification
*/

func ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Validating token...")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token missing", http.StatusUnauthorized)
			return
		}
		receivedToken := strings.Split(authHeader, " ")[1] // FIXME: access it safely to avoid runtime errors

		// another way to extrating the token - From Cookie
		// cookie, _ := r.Cookie("session_token")

		// 💡 Validating the Token
		claims := &UserClaims{}
		tkn, err := jwt.ParseWithClaims(receivedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwt_secret, nil
		})

		expiresAt, _ := tkn.Claims.GetExpirationTime()
		log.Println("Token will expire at: ", expiresAt)
		if err != nil || !tkn.Valid {
			http.Error(w, fmt.Sprintf("Invalid Token %v", err), http.StatusUnauthorized)
			return
		}

		// Extracting user data from token
		log.Println(claims.Username, "is requesting some data from protected route")

		//❤️  we can also set the userclaims in context and pass it to the main handler
		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RecoveryHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
