package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"secret-keeper/pkg/api/server"
	"strings"
)

type CLI struct {
	cl server.SecretKeeperClient
}

var ErrExit = errors.New("exit")

func New(addr string) (*CLI, error) {
	c := &CLI{}
	if err := c.connect(addr); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("Connected")
	return c, nil
}

func (c *CLI) Close() error {
	// TODO: close connection
	return nil
}

const (
	get  = "GET ◀️"
	set  = "SET ▶️"
	del  = "DELETE 🗑"
	back = "BACK ⬅️"
	web  = "WEBISTE 🕸"
)

var minCharacters = 8

const exit = "EXIT 🚪"
const auth = "SIGN IN 👤"
const reg = "SIGN UP 🆕"

var ErrUnavailable = errors.New("service unavailable")

func (c *CLI) Start(ctx context.Context) error {
	var header metadata.MD // variable to store header and trailer
	md := metadata.New(map[string]string{})
	ctx = metadata.NewOutgoingContext(ctx, md)

	if err := c.authenticate(ctx, &header); err != nil {
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
	if err := c.operate(ctx, &header, token); err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to operate: %w", err)
	}

	return nil
}

func (c *CLI) operate(ctx context.Context, header *metadata.MD, token string) error {
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	sp := selection.New("Choose", []string{
		get,
		set,
		del,
		web,
		exit})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		choice, err := sp.RunPrompt()
		if err != nil {
			return fmt.Errorf("failed to run prompt: %w", err)
		}

		choice = strings.Trim(choice, "\n\r")

		switch choice {
		case exit:
			return nil
		case get:
			key, backToMenu, err := c.getOneFromList(ctx, header)
			if err != nil {
				return fmt.Errorf("failed to get from list: %w", err)
			}

			if backToMenu {
				continue
			}

			if r, err := c.cl.Get(ctx, &server.GetRequest{Key: strings.Trim(key, "\n\r")},
				grpc.Header(header)); err != nil {
				log.Printf("Failed to get: %v", err)
			} else {
				log.Printf("Value: %s", r.Value)
			}
		case set:
			keyInput := textinput.New("Key:")
			keyInput.Placeholder = fmt.Sprintf("name of your secret")
			key, err := keyInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			valueInput := textinput.New("Value:")
			valueInput.Placeholder = fmt.Sprintf("value of your secret")
			value, err := valueInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			_, err = c.cl.Set(ctx, &server.SetRequest{Key: strings.Trim(key, "\n\r"),
				Value: strings.Trim(value, "\n\r")}, grpc.Header(header))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					log.Printf("Failed to set: %v", err)
				} else {
					if st.Code() == codes.Unavailable {
						log.Printf("Failed to set: %v", ErrUnavailable)
					} else {
						log.Printf("Failed to set: %v", err)
					}
				}
			} else {
				log.Print("OK\n")
			}
		case del:
			key, backToMenu, err := c.getOneFromList(ctx, header)
			if err != nil {
				return fmt.Errorf("failed to get from list: %w", err)
			}

			if backToMenu {
				continue
			}

			if _, err := c.cl.Delete(ctx, &server.DeleteRequest{Key: strings.Trim(key, "\n\r")},
				grpc.Header(header)); err != nil {
				st, ok := status.FromError(err)
				if !ok {
					log.Printf("Failed to delete: %v", err)
				} else {
					if st.Code() == codes.Unavailable {
						log.Printf("Failed to delete: %v", ErrUnavailable)
					} else if st.Code() == codes.NotFound {
						log.Printf("Value not found: %v", key)
					} else {
						log.Printf("Failed to delete: %v", err)
					}
				}
			} else {
				log.Printf("Deleted: %s", key)
			}
		}
	}
}

func (c *CLI) getOneFromList(ctx context.Context, header *metadata.MD) (string, bool, error) {
	getAllNames, err := c.cl.GetAllNames(ctx, &server.GetAllNamesRequest{}, grpc.Header(header))
	if err != nil {
		return "", false, fmt.Errorf("failed to get all: %w", err)
	}

	var names = make([]string, 0, len(getAllNames.Vars)+1)
	names = append(names, back)
	names = append(names, getAllNames.Vars...)

	msg := "Choose secret"
	if len(names) == 1 {
		msg = "No secrets"
	}

	getAllInput := selection.New(msg, names)
	getAllInput.PageSize = 10

	choice, err := getAllInput.RunPrompt()
	if err != nil {
		return "", false, fmt.Errorf("failed to run prompt: %w", err)
	}

	choice = strings.Trim(choice, "\n\r")

	if choice == back {
		return "", true, nil
	}
	return choice, false, nil
}

func (c *CLI) authenticate(ctx context.Context, header *metadata.MD) error {
	authInput := selection.New("Choose", []string{
		auth,
		reg,
		exit})
	authInput.PageSize = 3

	passInput := textinput.New("Passphrase:")
	passInput.Placeholder = fmt.Sprintf("more than %d characters", minCharacters)
	passInput.Validate = func(s string) error {
		if len(s) < minCharacters {
			return fmt.Errorf("at least %d more characters", minCharacters-len(s))
		}

		return nil
	}
	passInput.Hidden = true

	usernameInput := textinput.New("Username:")
	usernameInput.Placeholder = fmt.Sprintf("nickname")

	for {

		choice, err := authInput.RunPrompt()
		if err != nil {
			return fmt.Errorf("failed to run prompt: %w", err)
		}

		cmd := strings.Trim(choice, "\n\r")
		switch cmd {
		case auth:
			username, err := usernameInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			password, err := passInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			password = strings.Trim(password, "\n\r")

			_, err = c.cl.Auth(ctx, &server.AuthRequest{Username: username, Password: password}, grpc.Header(header))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					return fmt.Errorf("failed to auth: %w", err)
				}

				if st.Code() == codes.Unavailable {
					return ErrUnavailable
				}

				if st.Code() == codes.NotFound {
					fmt.Printf("Wrong password or username! \n")
					continue
				}

				return fmt.Errorf("failed to auth: %w", err)
			}
			return nil
		case reg:
			username, err := usernameInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			password, err := passInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			password = strings.Trim(password, "\n\r")

			_, err = c.cl.Register(ctx, &server.RegisterRequest{Username: username, Password: password}, grpc.Header(header))
			if err != nil {
				return fmt.Errorf("failed to auth: %w", err)
			}
			return nil
		case exit:
			return ErrExit
		}
	}
}

func (c *CLI) connect(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.cl = server.NewSecretKeeperClient(conn)
	return nil
}
