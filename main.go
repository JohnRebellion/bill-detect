package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var base = "/home/johnn/mnv2"

func run(args ...string) (bytes.Buffer, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = base
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	fmt.Println("Result: " + out.String())
	return out, err
}

func main() {
	fmt.Println("Hello World")
	app := fiber.New()
	app.Use(logger.New())
	app.Post("/", func(c *fiber.Ctx) error {
		prediction := []byte("Please try again")
		file, err := c.FormFile("image")

		// Check for errors:
		if err == nil {
			err = c.SaveFile(file, fmt.Sprintf("%s/%s", base, "image.jpg"))

			if err == nil {
				result, err := run("python3.11", "detect.py", fmt.Sprintf("%s/%s", base, "image.jpg"))
				if err == nil {
					if result.Bytes() != nil {
						prediction = result.Bytes()
						log.Println(string(prediction))
					}
				}
			}
		}

		return c.Send(prediction)
	})

	log.Fatal(app.Listen(":8000"))
}
