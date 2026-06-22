package main

import (
	"fmt"
	"hrms/api/routes"
	"hrms/internals/config"
	"hrms/internals/connection"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {

	conn, err := connection.InitDB()
	if err != nil {
		fmt.Println(err, conn)
		return
	}

	f := fiber.New()

	f.Use(recoverer.New(recoverer.Config{
		EnableStackTrace: true,
	}))

	routes.Routes(conn, f)

	fmt.Println("Server is running")
	if err := f.Listen(config.GetPort()); err != nil {
		fmt.Println("Server error:", err)
		return
	}

}
