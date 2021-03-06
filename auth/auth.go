/*
*
* Project printServer, an api rest responsible for printing to a Zebra thermal printer.
* The printServer will receive a cryptographic POST containing
* the Zpl content so that the printer can print.
*
* @package     authentication
* @author      @jeffotoni
* @size        11/08/2017
*
 */

package auth

import (
	"crypto/rsa"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeffotoni/jwt/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	pathPrivate = "./private.rsa"
	pathPublic  = "./public.rsa.pub"

	ProjectTitle = "jwt example"

	ExpirationHours = 24 // Hours
	DayExpiration   = 10 // Days

	//
	// md5(userR)
	//
	UserR = "e8961ed73a43314af2cce47a4f5b65db"

	//
	// md5(PassR)
	//
	PassR = "da124a000515052c1e5bc6e4748c77e5"
)

//
// Structure of our server configurations
//
type JsonMsg struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

//
// jwt init
//
func init() {

	//
	//
	//
	privateByte, err := ioutil.ReadFile(pathPrivate)

	if err != nil {

		WriteJson("error", "Private key not found!")
		return
	}

	//
	//
	//
	publicByte, errx := ioutil.ReadFile(pathPublic)

	if errx != nil {

		WriteJson("error", "Public key not found!")
		return
	}

	//
	//
	//
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateByte)

	if err != nil {

		WriteJson("error", "Could not parse privatekey!")
		return
	}

	//
	//
	//
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicByte)

	if err != nil {

		WriteJson("error", "ould not parse publickey!")
		return
	}
}

//
// jwt GenerateJWT
//
func GenerateJWT(model models.User) (string, string) {

	//
	// Generating date validation to return to the user
	//
	Expires := time.Now().Add(time.Hour * 24 * 10).Unix()

	//
	// convert int64
	//
	ExpiresInt64, _ := strconv.ParseInt(fmt.Sprintf("%v", Expires), 10, 64)

	//
	// convert time unix to Date RFC
	//
	ExpiresDateAll := time.Unix(ExpiresInt64, 0)

	//
	// Date
	//
	ExpiresDate := ExpiresDateAll.Format("2009-09-02")

	//
	// claims Token data, the header
	//
	claims := models.Claim{

		User: model.Login,
		StandardClaims: jwt.StandardClaims{

			//
			// Expires in 24 hours * 10 days
			//
			ExpiresAt: Expires,
			Issuer:    ProjectTitle,
		},
	}

	//
	// Generating token
	//
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//
	// Transforming into string
	//
	tokenString, err := token.SignedString(privateKey)

	if err != nil {

		return "Could not sign the token!", "2006-01-02"
	}

	//
	// return token string
	//
	return tokenString, ExpiresDate
}

//
// base64 (md5(key))
//
// login e password default in base 64
// curl -X POST -H "Content-Type: application/json"
// -H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
// "https://localhost:9001/token"
//
func AuthBasic(w http.ResponseWriter, r *http.Request) {

	//
	// Authorization Basic base64 Encode
	//
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {

		//
		//
		//
		HttpWriteJson(w, "error", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	//
	//
	//
	tokenBase64 := strings.Trim(auth[1], " ")

	//
	//
	//
	tokenBase64 = strings.TrimSpace(tokenBase64)

	//
	// token 64
	//
	authToken64 := strings.SplitN(tokenBase64, ":", 2)

	if len(authToken64) != 2 || authToken64[0] == "" || authToken64[1] == "" {

		HttpWriteJson(w, "error", "token base 64 invalid!", http.StatusUnauthorized)
		return
	}

	//
	//
	//
	tokenUserEnc := authToken64[0]

	//
	//
	//
	keyUserEnc := authToken64[1]

	//
	// User, Login byte
	//
	tokenUserDecode, _ := b64.StdEncoding.DecodeString(tokenUserEnc)

	//
	// key user byte
	//
	keyUserDec, _ := b64.StdEncoding.DecodeString(keyUserEnc)

	//
	// User, Login string
	//
	tokenUserDecodeS := strings.TrimSpace(strings.Trim(string(tokenUserDecode), " "))

	//
	// key user, string
	//
	keyUserDecS := strings.TrimSpace(strings.Trim(string(keyUserDec), " "))

	//
	// Validate user and password in the database
	//
	if tokenUserDecodeS == UserR && keyUserDecS == PassR {

		var model models.User

		model.Login = tokenUserDecodeS
		//model.Password = keyUserDec
		model.Password = ""
		model.Role = "admin"

		token, expires := GenerateJWT(model)

		result := models.ResponseToken{token, expires}
		jsonResult, err := json.Marshal(result)

		if err != nil {

			// fmt.Fprintln(w, "Error generating json!")
			HttpWriteJson(w, "error", "json.Marshal error generating!", http.StatusUnauthorized)
			return
		}

		//
		//
		//
		w.WriteHeader(http.StatusOK)

		//
		//
		//
		w.Header().Set("Content-Type", "application/json")

		//
		//
		//
		w.Write(jsonResult)

		/**
		{
		  "Token": "39a3099b45634f6eb511991fddde83752_v2",
		  "Expires": "2026-09-14"
		}
		*/

	} else {

		stringErr := "Invalid User or Key!"

		//
		//
		//
		w.WriteHeader(http.StatusForbidden)

		//
		//
		//
		w.Header().Set("Content-Type", "application/json")

		//
		//
		//
		HttpWriteJson(w, "error", stringErr, http.StatusUnauthorized)
	}

	//HttpWriteJson(w, "success", http.StatusText(http.StatusOK), http.StatusOK)

	defer r.Body.Close()
}

//
// HandlerValidate
//
func HandlerValidate(w http.ResponseWriter, r *http.Request) bool {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {

		//http.Error(w, "authorization failed", http.StatusUnauthorized)
		HttpWriteJson(w, "error", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	token := strings.Trim(auth[1], " ")
	strings.TrimSpace(token)

	// star
	parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

		return publicKey, nil

	})

	if err != nil || !parsedToken.Valid {

		//w.WriteHeader(http.StatusAccepted)
		//fmt.Fprintln(w, "Your token has expired!")
		HttpWriteJson(w, "error", "Your token has expired!", http.StatusAccepted)
		return false

	}

	claims, ok := parsedToken.Claims.(*models.Claim)

	if !ok || claims.User != UserR {

		//w.WriteHeader(http.StatusAccepted)
		//HttpWriteJson(w, "error", "There's something strange about your token!", http.StatusAccepted)
		fmt.Fprintln(w, "There's something strange about your token")
		return false
	}

	// fmt.Println("User: ", claims.User)

	//HttpWriteJson(w, "success", "Your token it's ok ["+claims.User+"]", http.StatusOK)
	//func2(w, r)
	return true
}

//
// Returns json without typing in http
//
func WriteJson(Status string, Msg string) {

	msgJsonStruct := &JsonMsg{Status, Msg}

	msgJson, errj := json.Marshal(msgJsonStruct)

	if errj != nil {

		fmt.Println(`{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	fmt.Println(msgJson)
}

//
// Returns json by typing on http
//
func HttpWriteJson(w http.ResponseWriter, Status string, Msg string, httpStatus int) {

	msgJsonStruct := &JsonMsg{Status, Msg}

	msgJson, errj := json.Marshal(msgJsonStruct)

	if errj != nil {

		fmt.Fprintln(w, `{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	w.WriteHeader(httpStatus)

	w.Header().Set("Content-Type", "application/json")

	w.Write(msgJson)
}
