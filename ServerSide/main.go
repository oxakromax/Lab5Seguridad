package main

import (
	"Lab4Client/dh"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

var Dh = new(dh.Dh)

func main() {
	Dh.Configure()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(Dh)
	})
	app.Post("/", func(ctx *fiber.Ctx) error {
		err := Dh.ReceiveJson(ctx.Body())
		if err != nil {
			return err
		}
		log.Println(Dh.GetFinalKey())
		return ctx.JSON(Dh)
	})
	app.Post("/confirm/", func(ctx *fiber.Ctx) error {
		log.Println("Son iguales:", string(ctx.Body()) == Dh.GetFinalKey())
		return ctx.SendString(Dh.GetFinalKey())
	})
	port, ok := os.LookupEnv("PORT") // For Heroku or Elastic
	if !ok {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
