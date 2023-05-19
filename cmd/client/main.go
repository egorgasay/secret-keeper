package main

import (
	"context"
	"log"
	"secret-keeper/internal/client/cli"
)

func main() {
	c := cli.CLI{}
	err := c.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
