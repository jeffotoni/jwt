/*
*
* Program to do a load test, and validate the response team of our api
*
* @package     main
* @author      @jeffotoni
* @size        29/07/2017
*
 */

// _[Rate limiting](http://en.wikipedia.org/wiki/Rate_limiting)_
// is an important mechanism for controlling resource
// utilization and maintaining quality of service. Go
// elegantly supports rate limiting with goroutines,
// channels, and [tickers](tickers).

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (

	//
	// base64 (md5(userR))
	//
	UserR = "ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI="

	//
	// base64 (md5(PassR))
	//
	PassR = "ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
)

type Ping struct {
	Msg string `json:"msg"`
}

type TokenStruct struct {
	Token string `json:"token"`
}

func ShootUrl(Url string, Token string) string {

	var ping = &Ping{}

	req, err := http.NewRequest("POST", Url, nil)

	req.Header.Set("X-Custom-Header", "valueHeader")

	req.Header.Set("Authorization", "Bearer "+Token)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		panic(err)
	}

	defer resp.Body.Close()

	//
	//
	//
	fmt.Println("Status: ", resp.Status)

	//
	//
	//
	if resp.Status == "200 OK" {

		// fmt.Println("response Status:", resp.Status)

		fmt.Println("response Headers:", resp.Header)

		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Println("response Body:", string(body))

		json.Unmarshal([]byte(string(body)), &ping)

		//
		//
		//
		msg2 := ping.Msg

		//
		//
		//
		ping.Msg = ""

		return string(msg2)

	} else {

		ping.Msg = ""
		return string("error")

	}
}

func GeToken(Url string, TokenAccess string, KeyAccess string) string {

	var token = &TokenStruct{}

	req, err := http.NewRequest("POST", Url, nil)

	//req.Header.Set("X-Custom-Header", "valueHeader")

	req.Header.Set("Authorization", "Basic "+TokenAccess+":"+KeyAccess)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Status", resp.Status)

	if resp.Status == "200 OK" {

		bodyToken, _ := ioutil.ReadAll(resp.Body)

		// fmt.Println("response Body:", string(body))

		json.Unmarshal([]byte(string(bodyToken)), &token)

		//
		//
		//
		tokenjson := token.Token

		//
		//
		//
		token.Token = ""

		return string(tokenjson)

	} else {

		token.Token = ""
		return string("error")

	}
}

func main() {

	endPoinToken := "http://localhost:9001/token"
	endPoint1 := "http://localhost:9001/ping"

	//
	// get token
	//

	TokenString := Shoot.GeToken(endPoinToken, UserR, PassR)

	TokenString = strings.TrimSpace(strings.Trim(TokenString, " "))

	fmt.Println("Token: ", TokenString)

	// First we'll look at basic rate limiting. Suppose
	// we want to limit our handling of incoming requests.
	// We'll serve these requests off a channel of the
	// same name.
	requests := make(chan int, 50)

	for i := 1; i <= 50; i++ {

		println("Loading requests: ", fmt.Sprintf("%d", i))
		time.Sleep(time.Millisecond * 40)
		requests <- i
	}

	close(requests)

	// This `limiter` channel will receive a value
	// every 100 or 300 milliseconds. This is the regulator in
	// our rate limiting scheme.
	limiter := time.Tick(time.Millisecond * 35)

	// time start
	//
	//
	time1 := time.Now()

	// By blocking on a receive from the `limiter` channel
	// before serving each request, we limit ourselves to
	// 1 request every 200 milliseconds.
	for req := range requests {

		<-limiter

		msg := Shoot.ShootUrl(endPoint1, TokenString)
		fmt.Println("request: ", req, "msg: ", msg)

		if req == 200 {

			fmt.Println("pause 2 segs")
			time.Sleep(time.Second * 2)
		}
	}

	time2 := time.Now()
	diff := time2.Sub(time1)

	fmt.Println(diff)
	fmt.Println("Enter enter to finish")

	var input string
	fmt.Scanln(&input)

}
