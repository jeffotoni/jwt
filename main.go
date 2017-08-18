/*
*
*
* @package     main
* @author      @jeffotoni
* @size        18/08/2017
*
 */

package main

//
//
//
import (
	"fmt"
	"github.com/didip/tollbooth"
	jwt "github.com/jeffotoni/jwt/auth"
	"log"
	"net/http"
	"time"
)

//
//
//
var (
	confServer *http.Server
	ServerPort = "9001"
)

//
//
//
const (
	HttpHeaderTitle = `Jwt Example`
	HttpHeaderMsg   = `Good Server, thank you!`

	SizeByteAllowed = 1 << 24
	NewLimiter      = 100

	HandlerToken = "/token"
	HandlerHello = "/hello"
	HandlerPing  = "/ping"
)

//
// Type responsible for defining a function that returns boolean
//
type fn func(w http.ResponseWriter, r *http.Request) bool

//
// Function responsible for abstraction and receive the
// authentication function and the handler that will execute if it is true
//
func HandlerFuncAuth(auth fn, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if auth(w, r) {

			handler(w, r)
		}
	}
}

//
// Mounting the properties on the api screen
//
func ShowScreen() {

	//
	// Basic Authentication
	//
	Token := "localhost:" + ServerPort + "" + HandlerToken

	//
	//
	//
	Ping := "localhost:" + ServerPort + "" + HandlerPing

	//
	//
	//
	Hello := "localhost:" + ServerPort + "" + HandlerHello

	//
	//
	//
	sizeMb := (SizeByteAllowed / 1024) / 1024

	//
	// Showing on the screen
	//
	fmt.Println("Start port:", ServerPort)
	fmt.Println("Endpoints:")
	fmt.Println(Token)
	fmt.Println(Ping)
	fmt.Println(Hello)
	fmt.Println("Max bytes:", sizeMb, "Mb")
}

//
// Testing whether the service is online
//
func Ping(w http.ResponseWriter, r *http.Request) {

	//
	//
	//
	json := `{"msg":"pong"}`

	//
	//
	//
	pong := []byte(json)

	//
	//
	//
	w.Header().Set(HttpHeaderTitle, HttpHeaderMsg)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-x83838374774")

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.WriteHeader(http.StatusOK)

	//
	//
	//
	w.Write(pong)
}

//
//
//
func Hello(w http.ResponseWriter, req *http.Request) {

	json := `{"status":"ok","msg":"hello"}`

	//
	//
	//
	hello := []byte(json)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-x83838374774")

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.WriteHeader(http.StatusOK)

	//
	//
	//
	w.Write(hello)
}

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
