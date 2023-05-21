package main

import (
	"context"
	"log"
	"secret-keeper/internal/client/cli"
)

func main() {
	c, err := cli.New("127.0.0.1:8080")
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	err = c.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
