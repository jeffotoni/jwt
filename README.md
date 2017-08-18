# Jwt example

The purpose of this project is to make a simplified version using JWT for educational purposes.

We will create a simplistic version for understanding, we will generate a token and the expiration date and send in json format to the requester.

The requester would have to have a way to register the user before using our authenticated API.

This part of the registry we will leave last, first we will simulate with keys that we create in the server for tests purposes.

We will Generate the Token and with the same token we will validate our handlers to see if they can be executed.

# Packages

go get -u github.com/dgrijalva/jwt-go

# Generate the keys

```sh

$ openssl genrsa -out private.rsa 1024

$ openssl rsa -in private.rsa -pubout > public.rsa.pub

```

# Install

$ go build main.go

$ sudo cp main /usr/bin/jwt

# Simulate 

$ go run simulate_ping.go


# Simulate Curl

$ curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU=" \
localhost:9001/token

$ curl -X POST -H "Content-Type: application/json" \
-H -H "Authorization: Bearer <TOKEN>" \
localhost:9001/ping

$ curl -X POST -H "Content-Type: application/json" \
-H -H "Authorization: Bearer <TOKEN>" \
localhost:9001/hello

# Main function

```go

//
// start
//
func main() {

	//
	//
	//
	ShowScreen()

	// Creating limiter for all handlers
	// or one for each handler. Your choice.
	// This limiter basically says: allow at most NewLimiter request per 1 second.
	limiter := tollbooth.NewLimiter(NewLimiter, time.Second)

	// Limit only GET and POST requests.
	limiter.Methods = []string{"GET", "POST"}

	//
	//
	//
	mux := http.NewServeMux()

	mux.Handle(HandlerPing, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(jwt.HandlerValidate, Ping)))

	mux.Handle(HandlerHello, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(jwt.HandlerValidate, Hello)))

	//
	// Off the default mux
	// Does not need authentication, only user key and token
	//
	mux.Handle(HandlerToken, tollbooth.LimitFuncHandler(limiter, jwt.AuthBasic))

	//
	//
	//
	confServer = &http.Server{

		Addr: ":" + ServerPort,

		Handler: mux,
		//ReadTimeout:    30 * time.Second,
		//WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 23, // Size accepted by package
	}

	log.Fatal(confServer.ListenAndServe())

}

```