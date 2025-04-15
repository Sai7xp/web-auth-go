## 1. Dealing with JSON

- JSON `Marshal` vs `UnMarshal`
- Be careful while marshalling `rune`
- Sending JSON in HTTP API response, converting go struct into JSON format and vice versa
- `http.ServeMux`, Creating a HTTP web server

## Authentication, Authorization

- **Authentication** - verifies who you are (Verifies that no-one is impersonating you)
- **Authorization** - Defines what we can do
- How can we add Authentication for APIs ?
  - One way is to do it via **Authorization** Header
  - Two common authorization schmes
  ```java
  Authorization: Basic QWxhZGRpbjpPcGVuU2VzYW1l
  Authorization: Bearer <token>
  ```
  - In **Basic** Authorization we will send "username:password" encoded in base64 format with every request
    - Never use Basic with http since we are sending username and password, only https is recommended
    - `req.BasicAuth()` in go will return the username and password, you don't have to handle the base64 decoding part
  - In **Bearer** we can send tokens like JWT.
  - To know more on Authorization Header - https://beeceptor.com/docs/concepts/authorization-header/

## Storing Passwords In DB
