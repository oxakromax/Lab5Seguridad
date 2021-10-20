package main

import (
	"Lab4/crypt"
	"Lab4/lab4"
	"Lab4/structs"
	"Lab4/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	lab4.DiffieHellmanSync() // LAB 4 Process
	Mensaje := File()        // Se lee el mensaje del archivo
	DES(Mensaje)             // Se encripta el mensaje con DES y se envia
	TripleDES(Mensaje)       // Lo mismo pero con 3DES
	AES(Mensaje)             // AES
}
func File() string {
	return string(utils.NewFileStore("Data/", ".txt").ReadRaw("mensajeentrada"))
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

func AES(Mensaje string) {
	i := utils.RandomInt(0, len(lab4.Dh.GetFinalKey())-32)
	key := []byte(lab4.Dh.GetFinalKey()[i : i+32])
	EncryptedText := crypt.AesEncrypt(key, []byte(Mensaje))
	MSG := &structs.MSG{
		Msg: EncryptedText,
		I:   i,
	}
	sendEncryptJson(MSG, "AES")
}

func TripleDES(Mensaje string) {
	i := utils.RandomInt(0, len(lab4.Dh.GetFinalKey())-24)
	key := []byte(lab4.Dh.GetFinalKey()[i : i+24])
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
	i := utils.RandomInt(0, len(lab4.Dh.GetFinalKey())-8)
	key := []byte(lab4.Dh.GetFinalKey()[i : i+8])
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
