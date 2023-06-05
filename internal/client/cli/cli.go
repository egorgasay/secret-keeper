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

// CLI is the command line interface
type CLI struct {
	logic *usecase.UseCase
}

// ErrExit is the exit error
var ErrExit = errors.New("exit")

// New creates a new CLI
func New(logic *usecase.UseCase) *CLI {
	return &CLI{logic: logic}
}

const (
	get  = "GET ◀️"
	set  = "SET ▶️"
	del  = "DELETE 🗑"
	back = "BACK ⬅️"
)

var minCharacters = 8

const (
	exit = "EXIT 🚪"
	auth = "SIGN IN 👤"
	reg  = "SIGN UP 🆕"
)

const (
	succeedAuth           = "Authenticated"
	chooseAction          = "Choose"
	noSecrets             = "No secrets"
	keyFieldName          = "Key: "
	keyFieldPlaceholder   = "name of your secret"
	valueFieldName        = "Value: "
	valueFieldPlaceholder = "your secret"
	PassphraseFieldName   = "Passphrase: "
	UserFieldName         = "Username: "
	UserFieldPlaceholder  = "nickname"
	NonUniqueUsername     = "Username already exists"
	InvalidCredentials    = "Invalid username or password"
)

// Start starts the CLI
func (c *CLI) Start(ctx context.Context) (err error) {
	ctx, err = c.authenticate(ctx)
	if err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	fmt.Println(succeedAuth)

	if err := c.operate(ctx); err != nil {
		if err == ErrExit {
			return nil
		}
		return fmt.Errorf("failed to operate: %w", err)
	}

	return nil
}

func trimNewlines(s string) string {
	return strings.Trim(s, "\n\r")
}

func (c *CLI) operate(ctx context.Context) error {
	sp := selection.New(chooseAction, []string{
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

		choice = trimNewlines(choice)

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

			secret, err := c.logic.GetSecret(ctx, trimNewlines(key))
			if err != nil {
				fmt.Printf("Failed to get: %v\n", err)
				continue
			}
			fmt.Printf("Secret: %s\n", secret)
		case set:
			keyInput := textinput.New(keyFieldName)
			keyInput.Placeholder = keyFieldPlaceholder
			key, err := keyInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			valueInput := textinput.New(valueFieldName)
			valueInput.Placeholder = valueFieldPlaceholder
			value, err := valueInput.RunPrompt()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			if err = c.logic.SetSecret(ctx, trimNewlines(key), trimNewlines(value)); err != nil {
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

			if err = c.logic.DeleteSecret(ctx, trimNewlines(key)); err != nil {
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

	msg := chooseAction
	if len(names) == 1 {
		msg = noSecrets
	}

	getAllInput := selection.New(msg, names)
	getAllInput.PageSize = 10

	choice, err := getAllInput.RunPrompt()
	if err != nil {
		return "", false, fmt.Errorf("failed to run prompt: %w", err)
	}

	choice = trimNewlines(choice)

	if choice == back {
		return "", true, nil
	}
	return choice, false, nil
}

func (c *CLI) authenticate(ctx context.Context) (context.Context, error) {
	authInput := selection.New(chooseAction, []string{
		auth,
		reg,
		exit})
	authInput.PageSize = 3

	passInput := textinput.New(PassphraseFieldName)
	passInput.Placeholder = fmt.Sprintf("more than %d characters", minCharacters)
	passInput.Validate = func(s string) error {
		if len(s) < minCharacters {
			return fmt.Errorf("at least %d more characters", minCharacters-len(s))
		}

		return nil
	}
	passInput.Hidden = true

	usernameInput := textinput.New(UserFieldName)
	usernameInput.Placeholder = UserFieldPlaceholder

	for {

		choice, err := authInput.RunPrompt()
		if err != nil {
			return ctx, fmt.Errorf("failed to run prompt: %w", err)
		}

		cmd := trimNewlines(choice)
		switch cmd {
		case auth:
			username, err := usernameInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read username: %w", err)
			}
			username = trimNewlines(username)

			password, err := passInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read password: %w", err)
			}

			password = trimNewlines(password)

			ctx, err = c.logic.Auth(ctx, username, password)
			if err != nil {
				if errors.Is(err, usecase.ErrInvalidPassword) {
					fmt.Println(InvalidCredentials)
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
			username = trimNewlines(username)

			password, err := passInput.RunPrompt()
			if err != nil {
				return ctx, fmt.Errorf("failed to read password: %w", err)
			}
			password = trimNewlines(password)

			ctx, err = c.logic.Register(ctx, username, password)
			if err != nil {
				if errors.Is(err, usecase.ErrUsernameExists) {
					fmt.Println(NonUniqueUsername)
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
