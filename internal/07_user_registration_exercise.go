package internal

import (
	"encoding/hex"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const (
	HTTP_PORT = ":9111"
)

// username : hashedpassword
var usersDb = make(map[string]string)

func RunUserRegExercise() {

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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Redirect(w, r, `/?errormsg=Invalid Credentials`, http.StatusSeeOther)
		return
	}
	_, isRegistered := usersDb[username]
	if isRegistered {
		http.Redirect(w, r, `/?errormsg=YOU ARE ALREADY REGISTERED!`, http.StatusSeeOther)
		return
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usersDb[username] = hex.EncodeToString(passwordHash)
	log.Println(passwordHash)

	// once the registration is successful, redirect user to login page
	http.Redirect(w, r, "/login-page", http.StatusSeeOther)
}

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
	_, isRegistered := usersDb[username]
	if !isRegistered {
		log.Println("User not registered")
		http.Redirect(w, r, `/?errormsg=Please Register to Login.`, http.StatusSeeOther)
		return
	}

	// if password matches then navigate to dashboard otherwise login page again
	existingPassword, _ := hex.DecodeString(usersDb[username])
	// log.Println(existingPassword)
	err := bcrypt.CompareHashAndPassword(existingPassword, []byte(password))
	if err != nil {
		log.Println("Password is incorrect", err)
		http.Redirect(w, r, `/login-page?errormsg=Invalid Credentials`, http.StatusSeeOther)
		return
	}

	log.Println("Yay! Login Successful")
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	userNameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login-page?errormsg=Not Logged In", http.StatusSeeOther)
		return
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>JWT Auth Exercise</title>
	</head>
	<body>
		<p> Dashboard </p>
		<p> Welcome,` + userNameCookie.Value + `</p>
	</body>
	</html>`

	w.Write([]byte(html))
}
