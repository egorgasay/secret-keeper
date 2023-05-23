package usecase

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"secret-keeper/pkg/api/server"
)

var ErrUnavailable = errors.New("service unavailable")
var ErrInvalidPassword = errors.New("Wrong password or username!")
var ErrUsernameExists = errors.New("username exists")
var ErrSecretNotFound = errors.New("secret exists")

type UseCase struct {
	cl     server.SecretKeeperClient
	header *metadata.MD
}

func New(addr string, header *metadata.MD) (*UseCase, error) {
	uc := &UseCase{header: header}
	if err := uc.connect(addr); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("Connected")
	return uc, nil
}

func (uc *UseCase) GetSecret(ctx context.Context, key string) (string, error) {
	r, err := uc.cl.Get(ctx, &server.GetRequest{Key: key}, grpc.Header(uc.header))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", fmt.Errorf("failed to get: %w", err)
		}
		if st.Code() == codes.Unavailable {
			return "", ErrUnavailable
		}
		if st.Code() == codes.NotFound {
			return "", ErrSecretNotFound
		}
		return "", err
	}
	return r.Value, nil
}

func (uc *UseCase) SetSecret(ctx context.Context, key, value string) error {
	_, err := uc.cl.Set(ctx, &server.SetRequest{Key: key, Value: value}, grpc.Header(uc.header))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to set: %w", err)
		} else {
			if st.Code() == codes.Unavailable {
				return fmt.Errorf("failed to set: %w", ErrUnavailable)
			}
			return fmt.Errorf("failed to set: %w", err)
		}
	}
	return nil
}

func (uc *UseCase) DeleteSecret(ctx context.Context, key string) error {
	_, err := uc.cl.Delete(ctx, &server.DeleteRequest{Key: key}, grpc.Header(uc.header))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to delete: %w", err)
		}
		if st.Code() == codes.Unavailable {
			return fmt.Errorf("failed to delete: %w", ErrUnavailable)
		} else if st.Code() == codes.NotFound {
			return fmt.Errorf("key not found: %v", key)
		}
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (uc *UseCase) GetAllNames(ctx context.Context) ([]string, error) {
	getAllNames, err := uc.cl.GetAllNames(ctx, &server.GetAllNamesRequest{}, grpc.Header(uc.header))
	if err != nil {
		return nil, fmt.Errorf("failed to get all: %w", err)
	}

	return getAllNames.Vars, nil
}

func (uc *UseCase) Auth(ctx context.Context, username, password string) (context.Context, error) {
	_, err := uc.cl.Auth(ctx, &server.AuthRequest{Username: username, Password: password}, grpc.Header(uc.header))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return ctx, fmt.Errorf("failed to auth: %w", err)
		}

		if st.Code() == codes.Unavailable {
			return ctx, ErrUnavailable
		}

		if st.Code() == codes.NotFound {
			return ctx, ErrInvalidPassword
		}

		return ctx, fmt.Errorf("failed to auth: %w", err)
	}

	return uc.addTokenToContext(ctx)
}

func (uc *UseCase) addTokenToContext(ctx context.Context) (context.Context, error) {
	tokens := uc.header.Get("token")
	if len(tokens) == 0 {
		return ctx, fmt.Errorf("failed to get token")
	}

	token := tokens[0]

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx, nil
}

func (uc *UseCase) Register(ctx context.Context, username, password string) (context.Context, error) {
	_, err := uc.cl.Register(ctx, &server.RegisterRequest{Username: username, Password: password}, grpc.Header(uc.header))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return ctx, fmt.Errorf("failed to register: %w", err)
		}

		if st.Code() == codes.Unavailable {
			return ctx, ErrUnavailable
		}

		if st.Code() == codes.AlreadyExists {
			return ctx, ErrUsernameExists
		}
		return ctx, fmt.Errorf("failed to auth: %w", err)
	}
	return uc.addTokenToContext(ctx)
}

func (uc *UseCase) connect(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	uc.cl = server.NewSecretKeeperClient(conn)
	return nil
}
