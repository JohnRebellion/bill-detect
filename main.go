package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func run(args ...string) (bytes.Buffer, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "/home/johnn/mnv2"
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
	app.Static("/", "./public")
	// app.Post("api/v1/detect", func(c *fiber.Ctx) error { return nil })
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			prediction := []byte("Please try again")
			if mt, msg, err = c.ReadMessage(); err != nil {
				// log.Println("read:", err)
				break
			}

			if strings.Contains(string(msg), "data:image/jpeg;base64,") ||
				strings.Contains(string(msg), "data:image/png;base64,") {
				result, err := run("python3.11", "detect.py", string(msg))

				if err == nil {
					if result.Bytes() != nil {
						prediction = result.Bytes()
						log.Println(string(prediction))
						// continue
					}
				}

				// log.Print(err)
			}

			if err = c.WriteMessage(mt, prediction); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	log.Fatal(app.Listen(":8000"))
}
