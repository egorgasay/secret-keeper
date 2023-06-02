package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	grpchandler "secret-keeper/internal/server/handler/grpc"
	"secret-keeper/internal/server/storage"
	"secret-keeper/internal/server/usecase"
	"secret-keeper/pkg"
	"secret-keeper/pkg/api/server"
	"testing"
)

func upServer() (stop func(), err error) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}
	defer store.Close()

	logic, err := usecase.New(store)
	if err != nil {
		log.Fatalf("failed to initialize logic: %v", err)
	}

	host := "127.0.0.1:8080"
	log.Printf("Server is running on grpc://%s\n", host)
	grpcServer := grpc.NewServer()
	ghandler := grpchandler.New(logic)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server.RegisterSecretKeeperServer(grpcServer, ghandler)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}
	}()

	ch := make(chan struct{})
	go func() {
		<-ch
		grpcServer.Stop()
		lis.Close()
	}()

	return func() {
		ch <- struct{}{}
	}, nil
}

func newEmptyContextWithMetadata() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
}

func TestUseCase_Auth(t *testing.T) {
	var header = metadata.MD{}

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:      newEmptyContextWithMetadata(),
				username: "XXXXX",
				password: "XXXXX",
			},
		},
		{
			name: "ok2",
			args: args{
				ctx:      newEmptyContextWithMetadata(),
				username: "admin",
				password: "test",
			},
		},
		{
			name: "noSuchUser",
			args: args{
				ctx:      newEmptyContextWithMetadata(),
				username: "noSuchUser",
				password: "XXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				newUUID, err := uuid.NewUUID()
				if err != nil {
					t.Errorf("Could not generate UUID %v", err)
					return
				}
				tt.args.username = newUUID.String()
				header.Set("key", tt.args.username)

				tt.args.ctx, err = uc.Register(tt.args.ctx, tt.args.username, tt.args.password)
				if err != nil && !errors.Is(err, ErrUsernameExists) {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			got, err := uc.Auth(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				md, ok := metadata.FromOutgoingContext(got)
				if !ok {
					t.Errorf("Auth() md == nil")
					return
				}
				if len(md) == 0 {
					t.Errorf("Auth() len(md) == 0, want != 0")
				}
				if len(md["token"]) == 0 {
					t.Errorf("Auth() len(md[\"token\"]) == 0, want != 0")
				}

				if len(md["token"][0]) == 0 {
					t.Errorf("Auth() len(md[\"token\"][0]) == 0, want != 0")
				}
			}
		})
	}
}

func TestUseCase_DeleteSecret(t *testing.T) {
	var header metadata.MD

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name     string
		args     args
		username string
		password string
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "TestUseCase_DeleteSecret1",
			},
			username: "XqfweXX",
			password: "wqfqwfq",
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				key: "TestUseCase_DeleteSecret2",
			},
			username: "XdqfXX",
			password: "XdqwdXXX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.args.ctx, err = uc.Register(tt.args.ctx, tt.username, tt.password)
				if errors.Is(err, ErrUsernameExists) {
					tt.args.ctx, err = uc.Auth(tt.args.ctx, tt.username, tt.password)
					if err != nil {
						t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
					}
				} else if err != nil {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				}
				err = uc.SetSecret(tt.args.ctx, tt.args.key, "XXXXX")
				if err != nil {
					t.Errorf("SetSecret() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if err = uc.DeleteSecret(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err = uc.GetSecret(tt.args.ctx, tt.args.key)
				if !errors.Is(err, ErrSecretNotFound) {
					t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestUseCase_GetAllNames(t *testing.T) {
	var header metadata.MD

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		username string
		want     []string
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			username: "GetAllNames1",
			want: []string{
				"TestUseCase_GetAllNames1",
				"TestUseCase_GetAllNames2",
			},
		},
		{
			name: "ok2",
			args: args{
				ctx: context.Background(),
			},
			username: "GetAllNames2",
			want: []string{
				"TestUseCase_GetAllNames3",
				"TestUseCase_GetAllNames4",
				"TestUseCase_GetAllNames5",
				"TestUseCase_GetAllNames6",
				"TestUseCase_GetAllNames7",
				"TestUseCase_GetAllNames8",
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
			},
			username: "GetAllNames3",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.args.ctx, err = uc.Register(tt.args.ctx, tt.username, tt.username)
				if errors.Is(err, ErrUsernameExists) {
					tt.args.ctx, err = uc.Auth(tt.args.ctx, tt.username, tt.username)
					if err != nil {
						t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
				for _, name := range tt.want {
					err = uc.SetSecret(tt.args.ctx, name, "XXXXX")
					if err != nil {
						t.Errorf("SetSecret() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
			}
			got, err := uc.GetAllNames(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !pkg.IsTheSameArray(got, tt.want) {
				t.Errorf("GetAllNames() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetSecret(t *testing.T) {
	var header metadata.MD

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name     string
		args     args
		username string
		want     string
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "TestUseCase_GetSecret",
			},
			want:    "XXXXX",
			wantErr: false,
		},
		{
			name: "notFound",
			args: args{
				ctx: context.Background(),
				key: "TestUseCase_GetSecret3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.args.ctx, err = uc.Register(tt.args.ctx, tt.username, tt.username)
				if errors.Is(err, ErrUsernameExists) {
					tt.args.ctx, err = uc.Auth(tt.args.ctx, tt.username, tt.username)
					if err != nil {
						t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
					}
				}

				err = uc.SetSecret(tt.args.ctx, tt.args.key, "XXXXX")
				if err != nil {
					t.Errorf("SetSecret() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			got, err := uc.GetSecret(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Register(t *testing.T) {
	var header metadata.MD

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name     string
		args     args
		username string
		want     string
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				ctx:      context.Background(),
				username: "TestUseCase_Register",
				password: "XXXXXXXXXXXXXXXXXXXX",
			},
			want:    "TestUseCase_Register",
			wantErr: false,
		},
		{
			name: "ok2",
			args: args{
				ctx:      context.Background(),
				username: "TestUseCase_Register2",
				password: "XXXXXXXXXXXXXXXXXXXX",
			},
			want:    "TestUseCase_Register2",
			wantErr: false,
		},
		{
			name: "notFound",
			args: args{
				ctx:      context.Background(),
				username: "TestUseCase_Register3",
				password: "XXXXXXXXXXXXXXXXXXXX",
			},
			want:    "TestUseCase_Register2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				newUUID, err := uuid.NewUUID()
				if err != nil {
					t.Errorf("Could not generate UUID %v", err)
					return
				}
				tt.args.username = newUUID.String()
			} else {
				uc.Register(tt.args.ctx, tt.args.username, tt.args.password)
			}
			ctx, err := uc.Register(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				md, ok := metadata.FromOutgoingContext(ctx)
				if !ok {
					t.Errorf("Could not get metadata")
					return
				}

				if len(md["token"]) == 0 || md["token"][0] == "" {
					t.Errorf("Could not get token")
				}
			}
		})
	}
}

func TestUseCase_SetSecret(t *testing.T) {
	var header metadata.MD

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				key:   "TestUseCase_SetSecret",
				value: "XXXXX",
			},
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				ctx:   context.Background(),
				key:   "TestUseCase_SetSecret",
				value: "XXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.args.ctx, err = uc.Register(tt.args.ctx, tt.args.key, tt.args.key)
				if errors.Is(err, ErrUsernameExists) {
					tt.args.ctx, err = uc.Auth(tt.args.ctx, tt.args.key, tt.args.key)
					if err != nil {
						t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
			}
			if err = uc.SetSecret(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := uc.GetSecret(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if got != tt.args.value {
					t.Errorf("GetSecret() got = %v, want %v", got, tt.args.value)
				}
			}
		})
	}
}

func TestUseCase_addTokenToContext(t *testing.T) {
	var header = metadata.MD{}

	uc, err := New("127.0.0.1:8080", &header)
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	stop, err := upServer()
	if err != nil {
		log.Fatalf("Could not start server %v", err)
	}
	defer stop()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			want:    "TestUseCase_addTokenToContext",
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				header.Set("token", tt.want)
			} else {
				header.Delete("token")
			}
			got, err := uc.addTokenToContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("addTokenToContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				md, ok := metadata.FromOutgoingContext(got)
				if !ok {
					t.Errorf("Could not get metadata")
					return
				}
				if len(md["token"]) == 0 || md["token"][0] == "" {
					t.Errorf("Could not get token")
					return
				}

				if md["token"][0] != tt.want {
					t.Errorf("addTokenToContext() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
