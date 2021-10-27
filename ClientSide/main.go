package main

import (
	// "Lab4/rsac"
	"Lab4/rsac"
	"Lab4/utils"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/openpgp/elgamal"
)

func getRSAPubkey() *rsa.PublicKey {
	resp, err := http.Get("http://localhost:8080/RSAgetpub")
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

func getGamalPubkey() *elgamal.PublicKey {
	resp, err := http.Get("http://localhost:8080/gamalgetpub")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var Pubk = new(elgamal.PublicKey)
	json.Unmarshal(body, Pubk)
	return Pubk
}

func gamalPostMessage(msg *utils.GamalMSG) {
	json, _ := json.Marshal(msg)
	responseBody := bytes.NewBuffer(json)
	resp, err := http.Post("http://localhost:8080/gamalDecrypt", "application/json", responseBody)
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

func RSApostMessage(message []byte) {
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
	GamalProcess(msg)
}

func GamalProcess(msg []byte) {
	var pubk = getGamalPubkey()
	c1, c2, err := elgamal.Encrypt(rand.Reader, pubk, msg)
	if err != nil {
		log.Fatalln(err)
	}
	Gamalmsg := &utils.GamalMSG{
		C1: c1,
		C2: c2,
	}
	gamalPostMessage(Gamalmsg)
}

func RSAProcess(msg []byte) {
	PubK := getRSAPubkey()
	EncryptedMsg := rsac.EncryptWithPublicKey(msg, PubK)
	RSApostMessage(EncryptedMsg)
}
