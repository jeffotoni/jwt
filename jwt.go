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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	pathPrivate = "./private.rsa"
	pathPublic  = "./public.rsa.pub"

	ProjectTitle = "printserver zebra"

	ExpirationHours = 24 // Hours
	DayExpiration   = 10 // Days

	UserR = "jeffotoni"
	PassR = "654323121"
)

//
//
//
const (
	SizeByteAllowed = 1 << 24

	HandlerToken = "/token"
	HandlerHello = "/hello"
	HandlerPing  = "/ping"
)

//
// User structure
//
type Login struct {

	//
	//
	//
	User string `json:"user"`

	//
	//
	//
	Password string `json:"password,omitempty"`

	//
	//
	//
	Role string `json:"role"`
}

//
// jwt
//
type Claim struct {

	//
	//
	//
	User string `json:"user"`

	//
	//
	//
	jwt.StandardClaims
}

//
// ResponseToken
//
type ResponseToken struct {

	//
	// token
	//
	Token string `json:"token"`

	Expires string `json:"expires"`
}

// This method Message is to return our messages
// in json, ie the client will
// receive messages in json format
type Message struct {
	Code int    `json:code`
	Msg  string `json:msg`
}

//
// Structure of our server configurations
//
type JsonMsg struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

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

// This method is a simplified abstraction
// so that we can send them to our client
// when making a request
func JsonMsg(codeInt int, msgText string) string {

	data := &Message{Code: codeInt, Msg: msgText}

	djson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(djson)
}

//
//
//
func check(e error) {

	if e != nil {

		panic(e)
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

	//
	//
	//
	mux := http.NewServeMux()

	mux.Handle(HandlerPing, HandlerFuncAuth(AutValidate, Ping))

	mux.Handle(HandlerHello, HandlerFuncAuth(AutValidate, Hello))

	//
	// Off the default mux
	// Does not need authentication, only user key and token
	//
	mux.Handle(HandlerToken, LoginBasic)

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
