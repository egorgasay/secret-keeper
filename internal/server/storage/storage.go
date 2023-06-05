package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"secret-keeper/pkg"
)

// Storage for data
type Storage struct {
	users  *itisadb.Index
	tokens *itisadb.Index
	logger pkg.Logger
}

// Config for storage
type Config struct {
	URI string
}

// New creates new storage
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

// ErrNotFound when value not found
var ErrNotFound = errors.New("not found")

// ErrUnavailable when storage is unavailable
var ErrUnavailable = errors.New("unavailable")

// ErrUnknown when unknown error
var ErrUnknown = errors.New("unknown")

// ErrAlreadyExists when something is already exists
var ErrAlreadyExists = errors.New("already exists")

// Get returns value by key
func (s *Storage) Get(ctx context.Context, username, key string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", s.handleIndexError(err)
	}

	v, err := index.Get(ctx, key)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}
		s.logger.Warn(fmt.Errorf("Storage.Get(): %w", err).Error())
		return "", ErrUnknown
	}
	return v, nil
}

func (s *Storage) handleIndexError(err error) error {
	if errors.Is(err, itisadb.ErrIndexNotFound) {
		// TODO: log error
		return ErrUnknown
	}
	if errors.Is(err, itisadb.ErrUnavailable) {
		return ErrUnavailable
	}
	s.logger.Warn(err.Error())
	return ErrUnknown
}

// Set adds k:v to storage
func (s *Storage) Set(ctx context.Context, username, key, value string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return s.handleIndexError(err)
	}

	err = index.Set(ctx, key, value, false)
	if err != nil {
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}

		s.logger.Warn(err.Error())
		return ErrUnknown
	}
	return nil
}

// AddToken adds token to storage
func (s *Storage) AddToken(ctx context.Context, token string, username string) error {
	err := s.tokens.Set(ctx, token, username, false)
	if err != nil {
		if errors.Is(err, itisadb.ErrUniqueConstraint) {
			return ErrAlreadyExists
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}
		s.logger.Warn(fmt.Errorf("Storage.AddToken(): %w", err).Error())
		return ErrUnknown
	}
	return nil
}

// GetUsername returns username of token
func (s *Storage) GetUsername(ctx context.Context, token string) (string, error) {
	username, err := s.tokens.Get(ctx, token)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}

		s.logger.Warn(fmt.Errorf("Storage.GetUsername(): %w", err).Error())
		return "", ErrUnknown
	}
	return username, nil
}

// AddUser adds user to storage
func (s *Storage) AddUser(ctx context.Context, username string, password string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return s.handleIndexError(err)
	}

	err = index.Set(ctx, "password", password, true)
	if err != nil {
		if errors.Is(err, itisadb.ErrUniqueConstraint) {
			return ErrAlreadyExists
		}

		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}

		s.logger.Warn(fmt.Errorf("Storage.AddUser(): %w", err).Error())
		return ErrUnknown
	}
	return nil
}

// GetPassword returns password of user
func (s *Storage) GetPassword(ctx context.Context, username string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", s.handleIndexError(err)
	}

	val, err := index.Get(ctx, "password")
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}

		s.logger.Warn(fmt.Errorf("Storage.GetPassword(): %w", err).Error())
		return "", ErrUnknown
	}
	return val, nil
}

// GetAllNames returns all names of user
func (s *Storage) GetAllNames(ctx context.Context, username string) ([]string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return nil, err
	}

	keyValues, err := index.GetIndex(ctx)
	if err != nil {
		return nil, s.handleIndexError(err)
	}

	var names []string

	delete(keyValues, "password")
	delete(keyValues, username)
	for key := range keyValues {
		names = append(names, key)
	}

	return names, nil
}

// Delete deletes key from index
func (s *Storage) Delete(ctx context.Context, username string, key string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return s.handleIndexError(err)
	}

	err = index.DeleteAttr(ctx, key)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}
		s.logger.Warn(fmt.Errorf("Storage.Delete(): %w", err).Error())

		return err
	}

	return nil
}
