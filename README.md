# Go x Echo API Boilerplate for me
Golang x Echoで認証をもつAPIを実装する際のボイラープレート
以下の機能を実装しています。


## Features
- Golang v 1.23.4
- Google Cloud Firestore
- Hashed password/Verify password
- Docker for Google Cloud Run
- endpoints: Signup/Signin
- endpoints: public_health, private_health

## Envs
<!-- Output level of the logger [trace, debug, info, warn, error] -->
LOG_LEVEL=debug
<!-- for test -->
ISTEST=test
<!-- Port to listen on -->
PORT=8080
<!-- Requested by the frontend
/ is not needed at the end -->
FRONTEND_URL=http://localhost:3000
<!-- Generate JWT Token by this secret -->
JWT_SECRET=secret
<!-- Google cloud project id -->
PROJECT_ID=xxx
