package sessionbasedauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	HTTP_PORT = ":9112"
	// hmac for creating a session token
	HMAC_SECRET = "hmac-secret-value"
)

// username : hashedpassword
var usersDB = make(map[string]string)

// uniqueSessionId : username
var sessionsDB = make(map[string]string)

func RunSessionBasedAuth() {

	mux := http.ServeMux{}

	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("POST /register", registerHandler)

	mux.HandleFunc("GET /login-page", loginPageHandler)
	mux.HandleFunc("POST /login", loginHandlerFunc)

	mux.HandleFunc("GET /dashboard", dashboardHandler) // for logged in users

	log.Printf("listening on port %s\n", HTTP_PORT)
	if err := http.ListenAndServe(HTTP_PORT, &mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// /register route will redirect to this page again, any error message will be set in "errormsg" query param
	queryParmError := r.URL.Query().Get("errormsg")

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>JWT Auth Exercise</title>
	</head>
	<body>
		<p> Register Form </p>
		<p> ` + queryParmError + `</p>
		<form action="/register" method="POST">
			<input name="username" placeholder="Username" />
			<input type="password" name="password" placeholder="Password" />
			<input type="submit" />
		</form>
		<br />
		<br />
		<a href="/login-page"> already registered ? Login here </a>
	</body>
	</html>`

	w.Write([]byte(html))
}

// new user registration - grab the username & password and store in db
func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Redirect(w, r, `/?errormsg=Invalid Credentials`, http.StatusSeeOther)
		return
	}
	_, isRegistered := usersDB[username]
	if isRegistered {
		http.Redirect(w, r, `/?errormsg=YOU ARE ALREADY REGISTERED!`, http.StatusSeeOther)
		return
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usersDB[username] = hex.EncodeToString(passwordHash)
	log.Println(passwordHash)

	// once the registration is successful, redirect user to login page
	http.Redirect(w, r, "/login-page", http.StatusSeeOther)
}

// returns html for login page
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	queryParmError := r.URL.Query().Get("errormsg")

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>JWT Auth Exercise</title>
	</head>
	<body>
		<p> LOGIN Form </p>
		<p> ` + queryParmError + `</p>
		<form action="/login" method="POST">
			<input name="username" placeholder="Username" />
			<input type="password" name="password" placeholder="Password" />
			<input type="submit" />
		</form>
		<a href="/"> NOT registered ? </a>
	</body>
	</html>`

	w.Write([]byte(html))
}

func loginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Redirect(w, r, `/login-page?errormsg=Username and Password are required`, http.StatusSeeOther)
		return
	}
	_, isRegistered := usersDB[username]
	if !isRegistered {
		log.Println("User not registered")
		http.Redirect(w, r, `/?errormsg=Please Register to Login.`, http.StatusSeeOther)
		return
	}

	// if password matches then navigate to dashboard otherwise login page again
	existingPassword, _ := hex.DecodeString(usersDB[username])
	// log.Println(existingPassword)
	err := bcrypt.CompareHashAndPassword(existingPassword, []byte(password))
	if err != nil {
		log.Println("Password is incorrect", err)
		http.Redirect(w, r, `/login-page?errormsg=Invalid Credentials`, http.StatusSeeOther)
		return
	}

	/*
		Correct password and username received from user. Now create a session and set that session token in Cookie
	*/

	// 🔒 Create a Session
	uniqueSessionId := username + time.Now().Format(time.RFC3339Nano) // usually in production session ids will be random, not something like this which is predictable
	sessionToken := createHmacToken(uniqueSessionId)
	sessionsDB[uniqueSessionId] = username // IN production we have to set the TTL for session
	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Second * 30), // session will be valid for 30 secs
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login-page?errormsg=Not Logged In", http.StatusSeeOther)
		return
	}
	// check if the session exists
	sessionId := verifyHmacToken(sessionCookie.Value)
	if sessionId == "" {
		http.Redirect(w, r, "/login-page?errormsg=Session token doesn't seems to be correct. Please Login again", http.StatusSeeOther)
		return
	}

	userName, ok := sessionsDB[sessionId]
	if !ok {
		http.Redirect(w, r, "/login-page?errormsg=Session Expired. Please Login again", http.StatusSeeOther)
		return
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Dashboard</title>
	</head>
	<body>
		<p> Dashboard </p>
		<p> Welcome,` + userName + " :  " + sessionId + `</p>
	</body>
	</html>`

	w.Write([]byte(html))
}

/*
	Creating HMAC tokens for Authentication
*/

// create a HMAC code (a signature) for the session id hmac = Hash(sessionId + secret)
func createHmacToken(sessionId string) string {
	h := hmac.New(sha256.New, []byte(HMAC_SECRET))
	h.Write([]byte(sessionId))
	signedHmacBytes := h.Sum(nil)

	// convert raw bytes to printable format - either hex or base64
	// let's go with hex format
	signatureAsHexString := hex.EncodeToString(signedHmacBytes)
	return signatureAsHexString + "|" + sessionId
}

// verifies the given sessionToken and returns sessionId
func verifyHmacToken(signedToken string) string {
	split := strings.Split(signedToken, "|")
	if len(split) != 2 {
		return ""
	}
	signature := split[0]
	sessionId := split[1]

	h := hmac.New(sha256.New, []byte(HMAC_SECRET))
	h.Write([]byte(sessionId))
	newSignature := hex.EncodeToString(h.Sum(nil))
	if newSignature != signature {
		return ""
	}
	return sessionId
}
