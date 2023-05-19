package storage

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
)

type Storage struct {
	users  *itisadb.Index
	tokens *itisadb.Index
}

type Config struct {
	URI string
}

func New(c Config) (*Storage, error) {
	db, err := itisadb.New(c.URI)
	if err != nil {
		return nil, err
	}

	tokens, err := db.Index(context.Background(), "tokens")
	if err != nil {
		return nil, err
	}

	users, err := db.Index(context.Background(), "users")
	if err != nil {
		return nil, err
	}

	return &Storage{
		users:  users,
		tokens: tokens,
	}, nil
}

func (s *Storage) Close() {
	// TODO: close connection
}

func (s *Storage) Get(ctx context.Context, username, key string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", fmt.Errorf("index: %w", err)
	}

	return index.Get(ctx, key)
}

func (s *Storage) Set(ctx context.Context, username, key, value string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return fmt.Errorf("index: %w", err)
	}

	return index.Set(ctx, key, value, false)
}

func (s *Storage) AddToken(ctx context.Context, token string, username string) error {
	return s.tokens.Set(ctx, token, username, true)
}

func (s *Storage) GetUsername(ctx context.Context, token string) (string, error) {
	return s.tokens.Get(ctx, token)
}

func (s *Storage) AddUser(ctx context.Context, username string, password string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return err
	}
	return index.Set(ctx, "password", password, false)
}

func (s *Storage) GetPassword(ctx context.Context, username string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", err
	}

	return index.Get(ctx, "password")
}
