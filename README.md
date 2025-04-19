## JSON on the Web

- JSON `Marshal` vs `UnMarshal`
- Be careful while marshalling `rune`
- Sending JSON in HTTP API response, converting go struct into JSON format and vice versa
- `http.ServeMux`, Creating a HTTP web server

## Authentication, Authorization

- **Authentication** - verifies who you are (verifies that no-one is impersonating you)
- **Authorization** - Defines what we can do
- How can we add Authentication for APIs ?
  - One way is to do it via **Authorization** Header
  - Two common authorization schmes
  ```java
  Authorization: Basic QWxhZGRpbjpPcGVuU2VzYW1l
  Authorization: Bearer <token>
  ```
  - In **Basic** Authorization we will send "username:password" encoded in base64 format with every request
    - **Never use Basic with http** since we are sending username and password, only **https** is recommended
    - `req.BasicAuth()` in go will return the username and password, you don't have to handle the base64 decoding part
  - In **Bearer** we can send tokens like JWT.
  - To know more on Authorization Header - https://beeceptor.com/docs/concepts/authorization-header/

## Hashing - Storing Passwords In DB

- **Hashing** algorithms will generate a **fixed array of bytes** for the given input.
  - Hashing algorithms like SHA-256 do the math on bits and produce raw binary output, not text.
  - **So if we do `string(bytes)` we will get gibberish or invisible chars. so we need to convert that bytes to printable format.**
  - We can either convert it to hexadecimal format or base64 format, which are printable formats.
  - Each hexadecimal is of **4 bits**, so for [32]byte hash we will get 64 length of hex string.
- Never store plain passwords in db
- Store password by Hashing them (hashing is irreversible). Even if the db gets leaked original password can't be recovered
- Go packages for hashing the data - `sha256`, `bcrypt`
- `shasum -a 256 go1.24.2.darwin-amd64.pkg`
- `echo -n "password" | shasum -a 256`

## Digital Signature

- It's like digital version of signing a document with pen - but mathematically secure.
- Use it for data **Integrity**(data is not altered) & **Authenticity**(msg came from expected person).
- **Sender**
  - Hash the message
  - Sign the Hash with sender **private key**
  - Send {message, signature}
- **Receiver**
  - Decrypt the signature using sender **public key**
  - Generate the hash of message
  - Compare both the hash values
- go package `ECDSA`

## HMAC

- Hash-based Message Authentication Code
- `hmac = Hash(message + secret)`
- secret is shared between two parties who trust each other
- Sender sends {hmac,message}
- Receiver **generates** the hmac of received message using same secret key and compares new hmac signature with received hmac code
- to prevent **faked bearer tokens**, use this hmac cryptographic "signing"

## Bearer Tokens

- added in http specification with OAUTH2
- uses authorization header & keyword "Bearer"

## JWT - JSON Web Tokens

- {JWT standard fields - header}.{Your fields - payload}.{Signature}
- `base64UrlEncode(header) + "." + base64UrlEncode(payload) + "." + signature`
- base64 encoding is used because it will not generate periods(.) and we want our jwt token to be divided into 3 parts
- [jwt.io](https://jwt.io)
- anyone can read the payload or header since these are just base64 encoded strings
- But the signature is generated using `HMAC SHA512`. hmac code is generated by combining `payload + secret` so no body else can fake it
- We have to verify every jwt token that we receive on backend. by regenerating the hmac code again(payload + secret). If it matches we can ensure that token is issued by us.
- Embedding usecase - [check here](internal/05_jwt_auth_api.go)
- [Authenticating HTTP APIs with JWT](internal/05_jwt_auth_api.go)
- **We can write a common middleware to validate the token, and if we want to pass the jwt payload (Claims) that we got by parsing the token to the main handler, we can pass it using the req Context.**

### Storing JWT in Cookies

- We can store the token in browser Cookies if the APIs are being consumed by the client app running on a Web Browser.
  - Set Cookie during login `http.SetCookie(w, cookie)`
  - Cookies will be added to the request automatically so we can easily retrieve it using `req.Cookie("name")`
- Storing jwt in cookies is not ideal for mobile apps

## base64

- Encode binary data into text. we typically use this encoding for transferring binary data over text based protocol like HTTP
- binary data - images, other files
- Read more here - https://sumanth.netlify.app/blog/03-charsets-encodings/
- go package `base64`
  - base64.URLEncoding
  - base64.StdEncoding

## [💎 Password Authentication to Website - Exercise][def]

- when user hits the '/' root url, we will return the html for registration page
  - User can enter username and password and hit enter, and the action for form is specified as POST in register form `<form action="/register" method="POST">`
- `/register` is a POST call used for registering a new user.
  - Store the new username and hashed password in db(`map[string]string`) and **redirect user** to Login page
  - `http.Redirect(w, r, "/login-page", http.StatusSeeOther)`
- `/login-page` is a GET request which returns the html for login form
  - input username and password and hit submit. the login req will go to `/login` post method
  - If the username and password matches redirect user to dashboard
- **The Problem:** We need password everytime to allow user to access a protected route.

> HTTP is stateless so for every request we need password and username to authenticate the request. Instead of this we can implement either **Session-based** auth or **JWT token** based Auth. We will ask for password, username only one time i.e., at the time of login and from next time a session will be created

## Session Based Authentication

- First User register or sign up with username and password
- Next at the time of login user will input the username and password. we will verify them with the values stored in db and make the login request success if they are correct.
- **Creating a session:**
  - Create a unique sessionId for user
  - Generate a hmac code with that session id - hmac = (sessionId + secret_key)
  - **Store the sessionId:username in db**
  - Set the `sessionSignature|sessionId` as Cookie.
- **Creating a session:**
  - 🫆 How do we verify that session token is valid from the next request ?
  - Extract the sessionId from the cookie in request (`sessionSignature|sessionId`) strings.Split(|)
  - Create a new Signature(hmac) with sessionId and secret key.
  - Compare the signature received in the cookie and the new generated signature. if both match then check the db.
  - If the session exists, we consider that session as valid. and we will get the userId from sessionDB where we mapped `sessionId:username`
- **Note:** Instead of hmac signature and all, we can also create the session token as a random string. (HDFC session token is just a random string generated using crypt/rand package)

[def]: internal/SessionBasedAndJWTAuthenticationExercise/01-PasswordBasedAuth/password_based_auth.go
