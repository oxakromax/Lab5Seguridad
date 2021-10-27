package main

import (
	"Lab4/rsac"
	"Lab4/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AppendToFile(data, method string) {
	utils.NewFileStore("OutPut/", ".txt").Write("mensajerecibido", method+":\n"+data)
}

func main() {
	Priv, Pub := rsac.GenerateKeyPair(4096)
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	app.Post("/rsacDecrypt", func(ctx *fiber.Ctx) error {
		data := new(utils.MSG)
		ctx.BodyParser(data)
		DecryptedBytes := rsac.DecryptWithPrivateKey(data.Msg, Priv)
		AppendToFile(string(DecryptedBytes), "RSA")
		msg := string(DecryptedBytes)
		return ctx.SendString(msg)
	})

	app.Get("/getpub", func(ctx *fiber.Ctx) error {
		return ctx.JSON(Pub)
	})

	port, ok := os.LookupEnv("PORT") // For Heroku or Elastic
	if !ok {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
