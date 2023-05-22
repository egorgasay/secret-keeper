package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"secret-keeper/internal/server/storage"
)

type IUseCase interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Auth(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
	GetAllNames(ctx context.Context) ([]string, error)
	Delete(ctx context.Context, key string) error
}

var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidPassword = errors.New("invalid password")

type getOrCreateFromContext func(ctx context.Context) (bool, string, error)

type UseCase struct {
	storage          *storage.Storage
	getOrCreateToken getOrCreateFromContext
}

func New(storage *storage.Storage) (*UseCase, error) {
	return &UseCase{
		storage:          storage,
		getOrCreateToken: getOrCreateToken,
	}, nil
}

func (u *UseCase) GetAllNames(ctx context.Context) ([]string, error) {
	username, err := u.getUsernameFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("getFromContext: %w", err)
	}

	return u.storage.GetAllNames(ctx, username)
}

func (u *UseCase) Get(ctx context.Context, key string) (string, error) {
	username, err := u.getUsernameFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("getFromContext: %w", err)
	}

	val, err := u.storage.Get(ctx, username, key)
	if err != nil {
		return "", fmt.Errorf("get: %w", err)
	}
	return val, nil
}

func (u *UseCase) Set(ctx context.Context, key, value string) error {
	username, err := u.getUsernameFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getFromContext: %w", err)
	}

	return u.storage.Set(ctx, username, key, value)
}

func (u *UseCase) Register(ctx context.Context, username string, password string) (string, error) {
	ok, token, err := u.getOrCreateToken(ctx)
	if err != nil {
		return "", fmt.Errorf("getOrCreateToken: %w", err)
	}

	if !ok {
		err = u.storeToken(ctx, token, username)
		if err != nil {
			return token, fmt.Errorf("storeToken: %w", err)
		}
	}

	err = u.storage.AddUser(ctx, username, password)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (u *UseCase) Auth(ctx context.Context, username string, password string) (string, error) {
	ok, token, err := u.getOrCreateToken(ctx)
	if err != nil {
		return "", fmt.Errorf("getOrCreateToken: %w", err)
	}

	if !ok {
		_, err = token, u.storeToken(ctx, token, username)
		if err != nil {
			return "", fmt.Errorf("storeToken: %w", err)
		}
	}

	if ok, err = u.validateToken(ctx, token); err != nil {
		return "", fmt.Errorf("validateToken: %w", err)
	} else if !ok {
		return "", fmt.Errorf("validateToken: %w", ErrInvalidToken)
	}

	// check if password is correct
	var passwordFromDB string
	if passwordFromDB, err = u.storage.GetPassword(ctx, username); err != nil { // TODO: index.Cmp(attrName, strToCmp)
		return "", fmt.Errorf("GetPassword: %w", err)
	}

	if passwordFromDB != password {
		return "", ErrInvalidPassword
	}

	return token, nil
}

func (u *UseCase) Delete(ctx context.Context, key string) error {
	username, err := u.getUsernameFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getFromContext: %w", err)
	}

	return u.storage.Delete(ctx, username, key)
}

// validateToken validates token
func (u *UseCase) validateToken(ctx context.Context, token string) (ok bool, err error) {
	if _, err = u.storage.GetUsername(ctx, token); err != nil {
		// TODO: better error recognition?
		return false, err
	}
	return true, nil
}

// storeToken stores token
func (u *UseCase) storeToken(ctx context.Context, token, username string) error {
	fmt.Println(token, username)
	// create a header that the gateway will watch for
	header := metadata.Pairs("token", token)
	// send the header back to the gateway
	if err := grpc.SetHeader(ctx, header); err != nil {
		return err
	}

	return u.storage.AddToken(ctx, token, username)
}

// generateToken generates token
func generateToken() (string, error) {
	val, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

func getOrCreateToken(ctx context.Context) (ok bool, token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		token, err = generateToken()
		return false, token, err
	}

	values := md.Get("token")
	if len(values) == 0 {
		token, err = generateToken()
		return false, token, err
	}

	return true, values[0], nil
}

func (u *UseCase) getUsernameFromContext(ctx context.Context) (username string, err error) {
	ok, token, err := u.getOrCreateToken(ctx)
	if err != nil {
		return "", fmt.Errorf("getOrCreateToken: %w", err)
	}

	if !ok {
		return "", fmt.Errorf("getOrCreateToken: %w", ErrInvalidToken)
	}

	username, err = u.storage.GetUsername(ctx, token)
	if err != nil {
		return "", fmt.Errorf("getUsername: %w", err)
	}

	return username, nil
}
