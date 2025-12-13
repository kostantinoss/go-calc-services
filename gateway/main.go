package main

import (
	"log"
	"test_system/gateway/internal"
)

func main() {
	gateway := &internal.ApiGateway{}
	if err := gateway.Init(); err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}

	gateway.Start()
}
