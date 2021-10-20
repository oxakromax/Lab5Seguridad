package main

import (
	"Lab4/crypt"
	"Lab4/crypt/dh"
	"Lab4/goRead"
	"Lab4/structs"
	"Lab4/utils"
	"bytes"
	"encoding/json"
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

func File() string {
	return string(goRead.NewFileStore("Data/", ".txt").ReadRaw("mensajeentrada"))
}

func DiffieHellmanSync() {
	PostRequest()
	log.Println(Dh.GetFinalKey())
	confirm()
}

func sendEncryptJson(jdes *structs.MSG, route string) {
	marshal, err := json.Marshal(jdes)
	if err != nil {
		return
	}
	responseBody := bytes.NewBuffer(marshal)
	_, err = http.Post("http://localhost:8080/"+route, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func main() {
	DiffieHellmanSync()
	Mensaje := File()
	DES(Mensaje)
	TripleDES(Mensaje)
	AES(Mensaje)
}
func AES(Mensaje string) {
	i := utils.RandomInt(0, len(Dh.GetFinalKey())-32)
	key := []byte(Dh.GetFinalKey()[i : i+32])
	EncryptedText := crypt.AesEncrypt(key, []byte(Mensaje))
	MSG := &structs.MSG{
		Msg: EncryptedText,
		I:   i,
	}
	sendEncryptJson(MSG, "AES")
}

func TripleDES(Mensaje string) {
	i := utils.RandomInt(0, len(Dh.GetFinalKey())-24)
	key := []byte(Dh.GetFinalKey()[i : i+24])
	EncryptedText, err := crypt.TripleDesEncrypt([]byte(Mensaje), key)
	if err != nil {
		log.Fatal(err)
	}
	MSG := &structs.MSG{
		Msg: EncryptedText,
		I:   i,
	}
	sendEncryptJson(MSG, "3DES")
}

func DES(Mensaje string) {
	i := utils.RandomInt(0, len(Dh.GetFinalKey())-8)
	key := []byte(Dh.GetFinalKey()[i : i+8])
	EncryptedText, err := crypt.DesEncrypt([]byte(Mensaje), key)
	if err != nil {
		log.Fatal(err)
	}
	MSG := &structs.MSG{
		Msg: EncryptedText,
		I:   i,
	}
	sendEncryptJson(MSG, "DES")
}
