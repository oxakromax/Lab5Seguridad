package main

import (
	"Lab4/gamal"
	"Lab4/rsac"
	"Lab4/utils"
	"log"
	"math/big"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/openpgp/elgamal"
)

func AppendToFile(data, method string) {
	utils.NewFileStore("OutPut/", ".txt").Write("mensajerecibido", method+":\n"+data)
}

func main() {
	RSAPriv, RSAPub := rsac.GenerateKeyPair(4096)
	GamalPriv := &elgamal.PrivateKey{
		PublicKey: elgamal.PublicKey{
			G: gamal.FromHex(gamal.GeneratorHex),
			P: gamal.FromHex(gamal.PrimeHex),
		},
		X: gamal.FromHex("42"),
	}
	GamalPriv.Y = new(big.Int).Exp(GamalPriv.G, GamalPriv.X, GamalPriv.P)
	GamalPub := &GamalPriv.PublicKey
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
		DecryptedBytes := rsac.DecryptWithPrivateKey(data.Msg, RSAPriv)
		AppendToFile(string(DecryptedBytes), "RSA")
		msg := string(DecryptedBytes)
		return ctx.SendString(msg)
	})

	app.Post("/gamalDecrypt", func(c *fiber.Ctx) error {
		data := new(utils.GamalMSG)
		c.BodyParser(data)
		DecryptedBytes, err := elgamal.Decrypt(GamalPriv, data.C1, data.C2)
		if err != nil {
			log.Fatal(err)
		}
		AppendToFile(string(DecryptedBytes), "Gamal")
		msg := string(DecryptedBytes)
		return c.SendString(msg)
	})

	app.Get("/gamalgetpub", func(c *fiber.Ctx) error {
		return c.JSON(GamalPub)
	})

	app.Get("/RSAgetpub", func(ctx *fiber.Ctx) error {
		return ctx.JSON(RSAPub)
	})

	port, ok := os.LookupEnv("PORT") // For Heroku or Elastic
	if !ok {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
