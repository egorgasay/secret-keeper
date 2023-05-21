package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"secret-keeper/internal/server/storage"
	"secret-keeper/internal/server/usecase"
	"secret-keeper/pkg/api/server"
)

type Handler struct {
	logic usecase.IUseCase
	server.UnimplementedSecretKeeperServer
}

func New(logic usecase.IUseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Auth(ctx context.Context, req *server.AuthRequest) (*server.AuthResponse, error) {
	_, err := h.logic.Auth(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) || errors.Is(err, usecase.ErrInvalidPassword) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &server.AuthResponse{}, nil
}

func (h *Handler) Register(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	_, err := h.logic.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}
	return &server.RegisterResponse{}, nil // TODO: REMOVE TOKEN?
}

func (h *Handler) Get(ctx context.Context, req *server.GetRequest) (*server.GetResponse, error) {
	v, err := h.logic.Get(ctx, req.GetKey())
	if err != nil {
		if err == storage.ErrNotFound {

		}
		return nil, err
	}
	return &server.GetResponse{Value: v}, nil
}

func (h *Handler) Set(ctx context.Context, req *server.SetRequest) (*server.SetResponse, error) {
	err := h.logic.Set(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		return nil, err
	}
	return &server.SetResponse{}, nil
}

func (h *Handler) GetAllNames(ctx context.Context, _ *server.GetAllNamesRequest) (*server.GetAllNamesResponse, error) {
	keys, err := h.logic.GetAllNames(ctx)
	if err != nil {
		return nil, err
	}
	return &server.GetAllNamesResponse{Vars: keys}, nil
}

func (h *Handler) Delete(ctx context.Context, req *server.DeleteRequest) (*server.DeleteResponse, error) {
	err := h.logic.Delete(ctx, req.GetKey())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &server.DeleteResponse{}, nil
}
