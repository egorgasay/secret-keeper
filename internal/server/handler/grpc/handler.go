package grpchandler

import (
	"context"
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
	token, err := h.logic.Auth(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &server.AuthResponse{Token: token}, nil
}

func (h *Handler) Register(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	token, err := h.logic.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &server.RegisterResponse{Token: token}, nil // TODO: REMOVE TOKEN?
}

func (h *Handler) Get(ctx context.Context, req *server.GetRequest) (*server.GetResponse, error) {
	v, err := h.logic.Get(ctx, req.GetKey())
	if err != nil {
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
