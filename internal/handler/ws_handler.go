package handler

import (
	// 1. Import our local package
	"user-age-api/internal/websocket"

	// 2. Import external package with an ALIAS (fiberws)
	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServeWS(hub *websocket.Hub) fiber.Handler {
	// Use fiberws.New (from the alias)
	return fiberws.New(func(c *fiberws.Conn) {
		
		// Create the client using our local struct
		client := &websocket.Client{
			Hub:  hub,
			Conn: c,
			Send: make(chan []byte, 256), // Now accessing the Public field 'Send'
		}

		// Use the Public field 'Register'
		client.Hub.Register <- client

		go client.WritePump()
		client.ReadPump()
	})
}