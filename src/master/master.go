package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/url"
)

type Slave struct {
	ID string
	URL string
}

func ParseGreeting(input []byte) Greeting {
	var data Greeting
	err := json.Unmarshal(input, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	return data
}

func ReadRequestBody(_ http.ResponseWriter, request *http.Request)  {
//	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
//	defer request.Body.Close()
//	var greeting Greeting
//	_ = json.Unmarshal(POSTrequestBody, &greeting)
//	return greeting.text
}

func handler(_ http.ResponseWriter, request *http.Request) {
	POSTrequestBody, _ := ioutil.ReadAll(request.Body)
	fmt.Printf("%T\n", POSTrequestBody)
	fmt.Printf("%T\n", request.Body)
//	defer request.Body.Close()

	var slave Slave
	_ = json.Unmarshal(POSTrequestBody, &slave)
	fmt.Println("SLAVE ID: ", slave.ID)
	fmt.Println("URL: ", slave.URL)

	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	var slaveAddress string
	if slave.ID == "1" {
		slaveAddress = raspberryPiIP["1"]
	}  else if slave.ID == "2" {
		slaveAddress = raspberryPiIP["2"]
	}

	form := url.Values{}
	form.Set("url", slave.URL)
	http.PostForm(slaveAddress, form)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}
