package lab4

import (
	"Lab4/crypt/dh"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

var Dh = new(dh.Dh)

func PostRequest() {
	Dh.Configure()
	responseBody := bytes.NewBuffer(Dh.MakeJson())
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8080/", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	_ = Dh.ReceiveJson(body)
}

func confirm() {
	responseBody := bytes.NewBuffer([]byte(Dh.GetFinalKey()))
	resp, err := http.Post("http://localhost:8080/confirm", "text/plain", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Son iguales:", string(body) == Dh.GetFinalKey())
}
func DiffieHellmanSync() {
	PostRequest()
	log.Println(Dh.GetFinalKey())
	confirm()
}
