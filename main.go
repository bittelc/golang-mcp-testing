package main

import (
	"log"

	"github.com/localrivet/gomcp/client"
)

func main() {
	// Create a new client
	c, err := client.NewClient("my-client",
		client.WithProtocolVersion("2025-03-26"),
		client.WithProtocolNegotiation(true),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Call a tool on the MCP server
	result, err := c.CallTool("say_hello", map[string]interface{}{
		"name": "World",
	})
	if err != nil {
		log.Fatalf("Tool call failed: %v", err)
	}

	log.Printf("Result: %v", result)
}
