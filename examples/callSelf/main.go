package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DAPR_URL = "http://localhost:3500/v1.0/invoke/"
)

func main() {
	app := fiber.New()
	client := &http.Client{}

	app.Get("/call", func(c *fiber.Ctx) error {
		nParam := c.Query("n")
		log.Printf("n = %s", nParam)

		n, _ := strconv.Atoi(nParam)
		nParam = strconv.Itoa(n - 1)

		if n == 0 {
			return c.SendString("n = 0")
		}

		req, _ := http.NewRequest("GET", DAPR_URL+"callSelf/method/call?n="+nParam, nil)
		req.Header.Set("tcc-id", c.Get("tcc-id"))
		client.Do(req)

		return c.SendString("n = " + nParam)
	})

	app.Listen(":3501")
}
