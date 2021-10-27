package main

import (
	// "Lab4/rsac"
	"Lab4/rsac"
	"Lab4/utils"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)


func getPubkey() *rsa.PublicKey {
	resp, err := http.Get("http://localhost:8080/getpub")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var Pubk = new(rsa.PublicKey)
	json.Unmarshal(body, Pubk)
	return Pubk
}

func postMessage(message []byte) {
	json, _ := json.Marshal(utils.MSG{
		Msg: message,
	})
	responseBody := bytes.NewBuffer(json)
	resp, err := http.Post("http://localhost:8080/rsacDecrypt", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}

func main() {
	msg := utils.NewFileStore("Input/", ".txt").ReadRaw("mensajeentrada")
	RSAProcess(msg)
}

func RSAProcess(msg []byte) {
	PubK := getPubkey()
	EncryptedMsg := rsac.EncryptWithPublicKey(msg, PubK)
	postMessage(EncryptedMsg)
}

