package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"secret-keeper/internal/client/cli"
	"secret-keeper/internal/client/usecase"
)

var (
	Version   = "N/A"
	BuildTime = "N/A"
)

const startText = `
Build version: %s
Build date: %s

`

func main() {
	fmt.Printf(startText, Version, BuildTime)

	var header metadata.MD // variable to store header and trailer
	md := metadata.New(map[string]string{})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	uc, err := usecase.New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	c := cli.New(uc)

	err = c.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
