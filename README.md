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

## Storing Passwords In DB - Hashing

- Never store plain passwords
- Store password by Hashing them (hashing is irreversible). Even if the db gets leaked original password can't be recovered
- Go packages for hashing the data - `sha256`, `bcrypt`
- `shasum -a 256 go1.24.2.darwin-amd64.pkg`
- `echo -n "password" | shasum -a 256`

## Digital Signature

- It's like digital version of signing a document with pen - but mathematically secure.
- Use it for data **Integrity**(data is not altered) & **Authenticity**(msg came from expected person).
- **Sender**
  - Hash the message
  - Sign the Hash with sender private key
  - Send {message, signature}
- **Receiver**
  - Decrypt the signature using sender public key
  - Generate the hash of message
  - Compare both the hash values

## HMAC

- Hash-based Message Authentication Code
- hmac = Hash(message + secret)
- secret is shared between two parties who trust each other
- Sender sends {hmac,message}
- Receiver generates the hmac of received message using secret key and compares new hmac signature with received one

