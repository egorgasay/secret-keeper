package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"secret-keeper/pkg/api/server"
	"strings"
)

type CLI struct {
	cl server.SecretKeeperClient
}

var ErrExit = errors.New("exit")

func (c *CLI) authenticate(ctx context.Context, stdin *bufio.Reader, header *metadata.MD) error {
	for {
		fmt.Print("auth - to authenticate\n")
		fmt.Print("reg - to register\n")
		fmt.Print("exit - to exit\n")
		fmt.Print("Enter command:")

		cmd, err := stdin.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read command: %w", err)
		}

		cmd = strings.Trim(cmd, "\n\r")
		switch cmd {
		case "auth":
			fmt.Print("USERNAME: ")
			username, err := stdin.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			fmt.Print("PASSWORD: ")
			password, err := stdin.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			password = strings.Trim(password, "\n\r")

			_, err = c.cl.Auth(ctx, &server.AuthRequest{Username: username, Password: password}, grpc.Header(header))
			if err != nil {
				return fmt.Errorf("failed to auth: %w", err)
			}
			return nil
		case "reg":
			fmt.Print("USERNAME: ")
			username, err := stdin.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			fmt.Print("PASSWORD: ")
			password, err := stdin.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			password = strings.Trim(password, "\n\r")

			_, err = c.cl.Register(ctx, &server.RegisterRequest{Username: username, Password: password}, grpc.Header(header))
			if err != nil {
				return fmt.Errorf("failed to auth: %w", err)
			}
			return nil
		case "exit":
			return ErrExit
		}
	}
}

func (c *CLI) connect() error {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.cl = server.NewSecretKeeperClient(conn)
	return nil
}

func (c *CLI) Close() error {
	// TODO: close connection
	return nil
}

func (c *CLI) operate(ctx context.Context, stdin *bufio.Reader, header *metadata.MD, token string) error {
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Print("Enter command: ")
		cmd, err := stdin.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read command: %w", err)
		}

		cmd = strings.Trim(cmd, "\n\r")
		switch cmd {
		case "exit":
			return nil
		case "get":
			fmt.Print("KEY: ")
			key, err := stdin.ReadString('\n')
			if err != nil {
				log.Printf("Failed to read key: %v", err)
				continue
			}

			if get, err := c.cl.Get(ctx, &server.GetRequest{Key: strings.Trim(key, "\n\r")},
				grpc.Header(header)); err != nil {
				log.Printf("Failed to get: %v", err)
			} else {
				log.Printf("Value: %s", get.Value)
			}
		case "set":
			fmt.Print("KEY: ")
			key, err := stdin.ReadString('\n')
			if err != nil {
				log.Printf("Failed to read key: %v", err)
				continue
			}

			fmt.Print("VALUE: ")
			value, err := stdin.ReadString('\n')
			if err != nil {
				log.Printf("Failed to read value: %v", err)
				continue
			}

			_, err = c.cl.Set(ctx, &server.SetRequest{Key: strings.Trim(key, "\n\r"),
				Value: strings.Trim(value, "\n\r")}, grpc.Header(header))
			if err != nil {
				log.Printf("Failed to set: %v", err)
			} else {
				log.Print("OK\n")
			}
		}
	}
}

func (c *CLI) Start(ctx context.Context) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("Connected")

	var header metadata.MD // variable to store header and trailer
	md := metadata.New(map[string]string{})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stdin := bufio.NewReader(os.Stdin)

	if err := c.authenticate(ctx, stdin, &header); err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	fmt.Println("Authenticated")

	tokens := header.Get("token")
	if len(tokens) == 0 {
		return fmt.Errorf("failed to get token")
	}

	token := tokens[0]
	if err := c.operate(ctx, stdin, &header, token); err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to operate: %w", err)
	}

	return nil
}
