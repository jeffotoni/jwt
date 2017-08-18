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
	"crypto/rsa"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/didip/tollbooth"
	auth "github.com/jeffotoni/jwt/auth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//
//
//
const (
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
	Token := "localhost:9001/" + HandlerToken

	//
	//
	//
	Ping := "localhost:9001/" + HandlerPing

	//
	//
	//
	Hello := "localhost:9001/" + HandlerHello

	//
	//
	//
	sizeMb := (SizeByteAllowed / 1024) / 1024

	//
	// Showing on the screen
	//
	fmt.Println("Start port:", 9001)
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
	w.Write(json)
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

	mux.Handle(HandlerPing, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(AutValidate, Ping)))

	mux.Handle(HandlerHello, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(AutValidate, Hello)))

	//
	// Off the default mux
	// Does not need authentication, only user key and token
	//
	mux.Handle(HandlerToken, tollbooth.LimitFuncHandler(limiter, AuthBasic))

	//
	//
	//
	confServer = &http.Server{

		Addr: ":" + cfg.ServerPort,

		Handler: mux,
		//ReadTimeout:    30 * time.Second,
		//WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 23, // Size accepted by package
	}

	log.Fatal(confServer.ListenAndServe())

}
