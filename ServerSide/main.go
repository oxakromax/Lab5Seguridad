package main

import (
	"Lab4/crypt"
	"Lab4/crypt/dh"
	"Lab4/goRead"
	"Lab4/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

var Dh = new(dh.Dh)

func SaveFileText(text, method string) {
	file := goRead.NewFileStore("Data/", ".txt")
	file.Write("mensajerecibido", "Desencriptado por el metodo ("+method+"):\n"+text)
}

func main() {
	Dh.Configure()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "POST",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	app.Post("/", PostROOT)
	app.Post("/confirm/", POSTConfirm)

	app.Post("/DES/", PostDES)
	app.Post("/3DES/", Post3DES)
	app.Post("/AES/", PostAES)
	port, ok := os.LookupEnv("PORT") // For Heroku or Elastic
	if !ok {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}

func POSTConfirm(ctx *fiber.Ctx) error {
	log.Println("Son iguales:", string(ctx.Body()) == Dh.GetFinalKey())
	return ctx.SendString(Dh.GetFinalKey())
}

func PostROOT(ctx *fiber.Ctx) error {
	err := Dh.ReceiveJson(ctx.Body())
	if err != nil {
		return err
	}
	log.Println(Dh.GetFinalKey())
	return ctx.JSON(Dh)
}

func PostDES(ctx *fiber.Ctx) error {
	jDes := new(structs.MSG)
	err := ctx.BodyParser(jDes)
	if err != nil {
		return err
	}
	i := jDes.I
	destext, err := crypt.DesDecrypt(jDes.Msg, []byte(Dh.GetFinalKey()[i:i+8]))
	if err != nil {
		return err
	}
	SaveFileText(string(destext), "DES")
	return nil
}
func Post3DES(ctx *fiber.Ctx) error {
	jDes := new(structs.MSG)
	err := ctx.BodyParser(jDes)
	if err != nil {
		return err
	}
	i := jDes.I
	destext, err := crypt.TripleDesDecrypt(jDes.Msg, []byte(Dh.GetFinalKey()[i:i+24]))
	if err != nil {
		return err
	}
	SaveFileText(string(destext), "3DES")
	return nil
}
func PostAES(ctx *fiber.Ctx) error {
	jDes := new(structs.MSG)
	err := ctx.BodyParser(jDes)
	if err != nil {
		return err
	}
	i := jDes.I
	destext := crypt.AesDesEncrypt([]byte(Dh.GetFinalKey()[i:i+32]), jDes.Msg)
	if err != nil {
		return err
	}
	SaveFileText(string(destext), "AES")
	return nil
}
