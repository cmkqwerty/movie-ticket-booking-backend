package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listener address of the HTTP API server.")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/init", handleInit)

	app.Listen(*listenAddr)
}

func handleInit(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "success"})
}
