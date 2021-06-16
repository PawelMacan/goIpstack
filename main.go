package main

import (
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"goipstack/service"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app := fiber.New()
	app.Post("/api/geoip/:ip?", service.CreateIpGeoLocation)
	app.Get("/api/geoip/:id?", service.GetIpGeoLocation)
	app.Delete("/api/geoip/:id?", service.DeleteIpGeoLocation)

	app.Listen(os.Getenv("SERVER_PORT"))
}
