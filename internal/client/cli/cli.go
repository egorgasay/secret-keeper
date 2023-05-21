package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"log"
	"secret-keeper/internal/client/usecase"
	"strings"
)

type CLI struct {
	logic *usecase.UseCase
}

var ErrExit = errors.New("exit")

func New(logic *usecase.UseCase) *CLI {
	return &CLI{logic: logic}
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
)

var minCharacters = 8

const exit = "EXIT 🚪"
const auth = "SIGN IN 👤"
const reg = "SIGN UP 🆕"

func (c *CLI) Start(ctx context.Context) (err error) {
	ctx, err = c.authenticate(ctx)
	if err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	fmt.Println("Authenticated")

	if err := c.operate(ctx); err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to operate: %w", err)
	}

	return nil
}

func (c *CLI) operate(ctx context.Context) error {
	sp := selection.New("Choose", []string{
		get,
		set,
		del,
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
			key, backToMenu, err := c.getOneFromList(ctx)
			if err != nil {
				return fmt.Errorf("failed to get from list: %w", err)
			}

			if backToMenu {
				continue
			}

			secret, err := c.logic.GetSecret(ctx, strings.Trim(key, "\n\r"))
			if err != nil {
				fmt.Printf("Failed to get: %v\n", err)
				continue
			}
			fmt.Printf("Secret: %s\n", secret)
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

			if err = c.logic.SetSecret(ctx, strings.Trim(key, "\n\r"), strings.Trim(value, "\n\r")); err != nil {
				log.Println(err)
			} else {
				log.Print("OK\n")
			}
		case del:
			key, backToMenu, err := c.getOneFromList(ctx)
			if err != nil {
				return fmt.Errorf("failed to get from list: %w", err)
			}

			if backToMenu {
				continue
			}

			if err = c.logic.DeleteSecret(ctx, strings.Trim(key, "\n\r")); err != nil {
				log.Println(err)
			} else {
				log.Printf("Deleted: %s", key)
			}
		}
	}
}

func (c *CLI) getOneFromList(ctx context.Context) (string, bool, error) {
	keys, err := c.logic.GetAllNames(ctx)
	if err != nil {
		return "", false, err
	}

	var names = make([]string, 0, len(keys)+1)
	names = append(names, back)
	names = append(names, keys...)

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

func (c *CLI) authenticate(ctx context.Context) (context.Context, error) {
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
			return ctx, fmt.Errorf("failed to run prompt: %w", err)
		}

		cmd := strings.Trim(choice, "\n\r")
		switch cmd {
		case auth:
			username, err := usernameInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			password, err := passInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read password: %w", err)
			}

			password = strings.Trim(password, "\n\r")

			ctx, err = c.logic.Auth(ctx, username, password)
			if err != nil {
				if errors.Is(err, usecase.ErrInvalidPassword) {
					fmt.Println("Invalid password or username")
					continue
				}
				return ctx, fmt.Errorf("failed to auth: %w", err)
			}
			return ctx, nil

		case reg:
			username, err := usernameInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read username: %w", err)
			}
			username = strings.Trim(username, "\n\r")

			password, err := passInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read password: %w", err)
			}
			password = strings.Trim(password, "\n\r")

			ctx, err = c.logic.Register(ctx, username, password)
			if err != nil {
				if errors.Is(err, usecase.ErrUsernameExists) {
					fmt.Println("Username already exists")
					continue
				}
				return ctx, fmt.Errorf("failed to register: %w", err)
			}
			return ctx, nil
		case exit:
			return ctx, ErrExit
		}
	}
}
