package storage

import (
	"context"
	"errors"
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

var ErrNotFound = errors.New("not found")
var ErrUnavailable = errors.New("unavailable")
var ErrUnknown = errors.New("unknown")
var ErrAlreadyExists = errors.New("already exists")

func (s *Storage) Close() {
	// TODO: close connection
}

func (s *Storage) Get(ctx context.Context, username, key string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", handleIndexError(err)
	}

	v, err := index.Get(ctx, key)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}
		// TODO: log error
		return "", ErrUnknown
	}
	return v, nil
}

func handleIndexError(err error) error {
	if errors.Is(err, itisadb.ErrIndexNotFound) {
		// TODO: log error
		return ErrUnknown
	}
	if errors.Is(err, itisadb.ErrUnavailable) {
		return ErrUnavailable
	}
	// TODO: log error
	return ErrUnknown
}

func (s *Storage) Set(ctx context.Context, username, key, value string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return handleIndexError(err)
	}

	err = index.Set(ctx, key, value, false)
	if err != nil {
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}

		// TODO: log error
		return ErrUnknown
	}
	return nil
}

func (s *Storage) AddToken(ctx context.Context, token string, username string) error {
	err := s.tokens.Set(ctx, token, username, false)
	if err != nil {
		if errors.Is(err, itisadb.ErrUniqueConstraint) {
			return ErrAlreadyExists
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}
		// TODO: log error
		return ErrUnknown
	}
	return nil
}

func (s *Storage) GetUsername(ctx context.Context, token string) (string, error) {
	username, err := s.tokens.Get(ctx, token)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}

		// TODO: log error
		return "", ErrUnknown
	}
	return username, nil
}

func (s *Storage) AddUser(ctx context.Context, username string, password string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return handleIndexError(err)
	}

	err = index.Set(ctx, "password", password, true)
	if err != nil {
		if errors.Is(err, itisadb.ErrUniqueConstraint) {
			return ErrAlreadyExists
		}

		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}

		// TODO: log error
		return ErrUnknown
	}
	return nil
}

func (s *Storage) GetPassword(ctx context.Context, username string) (string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return "", handleIndexError(err)
	}

	val, err := index.Get(ctx, "password")
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return "", ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return "", ErrUnavailable
		}

		// TODO: log error
		return "", ErrUnknown
	}
	return val, nil
}

func (s *Storage) GetAllNames(ctx context.Context, username string) ([]string, error) {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return nil, err
	}

	keyValues, err := index.GetIndex(ctx)
	if err != nil {
		return nil, handleIndexError(err)
	}

	var names []string

	delete(keyValues, "password")
	delete(keyValues, username)
	for key := range keyValues {
		names = append(names, key)
	}

	return names, nil
}

func (s *Storage) Delete(ctx context.Context, username string, key string) error {
	index, err := s.users.Index(ctx, username)
	if err != nil {
		return handleIndexError(err)
	}

	err = index.DeleteAttr(ctx, key)
	if err != nil {
		if errors.Is(err, itisadb.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, itisadb.ErrUnavailable) {
			return ErrUnavailable
		}

		return err
	}

	return nil
}
