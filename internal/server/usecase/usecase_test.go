package usecase

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"secret-keeper/internal/server/storage"
	"secret-keeper/pkg"
	"testing"
)

func TestUseCase_Auth(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		token   string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      setToken(context.Background(), "XXX-16tok"),
				username: "XXX-16",
				password: "XXX-16",
			},
			token: "XXX-16tok",
		},

		{
			name: "success2",
			args: args{
				ctx:      setToken(context.Background(), "qwerty"),
				username: "admin3",
				password: "tgfqdf",
			},
			token: "qwerty",
		},

		{
			name: "error",
			args: args{
				ctx:      setToken(context.Background(), ""),
				username: "XX2424rt52dXXXX",
				password: "XwdefwegfXXXXX",
			},
			wantErr: true,
		},

		{
			name: "noToken",
			args: args{
				ctx:      context.Background(),
				username: "XX2424rt52dXXXX",
				password: "XwdefwegfXXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			var token string
			if !tt.wantErr {
				token, err = u.Register(tt.args.ctx, tt.args.username, tt.args.password)
				if err != nil && !errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				err := u.storage.AddToken(tt.args.ctx, token, tt.args.username)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := u.Auth(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != token {
				t.Errorf("Auth() got = %v, want %v", got, token)
			}
		})
	}
}

type streamerStub struct {
}

func (s streamerStub) Method() string {
	//TODO implement me
	panic("implement me")
}

func (s streamerStub) SetHeader(md metadata.MD) error {
	return nil
}

func (s streamerStub) SendHeader(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func (s streamerStub) SetTrailer(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func setToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx
}

func setHeader(ctx context.Context) context.Context {
	return grpc.NewContextWithServerTransportStream(ctx, &streamerStub{})
}

func TestUseCase_Delete(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name     string
		args     args
		token    string
		username string
		password string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_Delete-163s"),
				key: "XXX-16",
			},
			username: "qwdq",
			password: "XXXX",
			token:    "TestUseCase_Delete-163s",
		},
		{
			name: "err",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_Delete-164s"),
				key: "XXX-32",
			},
			username: "qds",
			password: "dwqd",
			token:    "TestUseCase_Delete-164s",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				err := u.storage.AddToken(tt.args.ctx, tt.token, tt.username)
				if err != nil {
					t.Error(err)
				}

				err = u.Set(tt.args.ctx, tt.args.key, tt.username)
				if err != nil {
					t.Error(err)
				}
			}

			if err = u.Delete(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, err = u.Get(tt.args.ctx, tt.args.key); !errors.Is(err, storage.ErrNotFound) {
					t.Fatalf("wantErr %v got %v", storage.ErrNotFound, err)
				}
			}
		})
	}
}

func TestUseCase_Get(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		token    string
		username string
		password string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_Get1"),
				key: "XXX-16",
			},
			username: "qwdq",
			password: "XXXX",
			token:    "TestUseCase_Get1",
			want:     "qwdq",
		},
		{
			name: "err",
			args: args{
				ctx: setToken(context.Background(), ""),
				key: "XXX-32",
			},
			username: "qds",
			password: "dwqd",
			token:    "dqwd-16",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				err := u.storage.AddToken(tt.args.ctx, tt.token, tt.username)
				if err != nil {
					t.Error(err)
				}

				err = u.Set(tt.args.ctx, tt.args.key, tt.want)
				if err != nil {
					t.Error(err)
				}
			}

			got, err := u.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetAllNames(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		token    string
		username string
		password string
		want     []string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_GetAllNames-1"),
			},
			username: "XXXX455",
			password: "XXXX456",
			token:    "TestUseCase_GetAllNames-1",
			want: []string{
				"333", "444", "XXX44",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_GetAllNames-2"),
			},
			username: "XXXX4355",
			password: "XXXX456",
			token:    "TestUseCase_GetAllNames-2",
			want:     []string{},
		},
		{
			name: "error",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_GetAllNames-3"),
			},
			username: "X244355",
			password: "XXXX456",
			token:    "TestUseCase_GetAllNames-3",
			want:     []string{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				err = u.storage.AddToken(tt.args.ctx, tt.token, tt.username)
				if err != nil {
					t.Error(err)
				}

				for _, name := range tt.want {
					err = u.Set(tt.args.ctx, name, tt.username)
					if err != nil {
						t.Error(err)
					}
				}
			}

			got, err := u.GetAllNames(tt.args.ctx)
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

func TestUseCase_Register(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		token   string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      setToken(context.Background(), "TestUseCase_Register-1"),
				username: "XXX-16",
				password: "XXX-16",
			},
			token: "TestUseCase_Register-1",
		},

		{
			name: "success2",
			args: args{
				ctx:      setToken(context.Background(), "TestUseCase_Register-2"),
				username: "admin3",
				password: "tgfqdf",
			},
			token: "TestUseCase_Register-2",
		},

		{
			name: "alreadyExists",
			args: args{
				ctx:      setToken(context.Background(), "TestUseCase_Register-3"),
				username: "admin3",
				password: "XXXXXX",
			},
			token:   "TestUseCase_Register-3",
			wantErr: true,
		},
		{
			name: "noToken",
			args: args{
				ctx:      context.Background(),
				username: "admin3",
				password: "XXXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			got, err := u.Register(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				if !errors.Is(err, storage.ErrAlreadyExists) {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !tt.wantErr {
				err = u.storage.AddToken(tt.args.ctx, tt.token, tt.args.username)
				if err != nil {
					t.Error(err)
				}

				token, err := u.Auth(tt.args.ctx, tt.args.username, tt.args.password)
				if err != nil {
					t.Error(err)
				}
				if got != tt.token {
					t.Errorf("Register() token = %v, want %v", token, tt.token)
				}
			}
		})
	}
}

func TestUseCase_Set(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name     string
		args     args
		username string
		password string
		token    string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:   setToken(context.Background(), "TestUseCase_Set1"),
				key:   "333",
				value: "333",
			},
			username: "qwe",
			password: "XXX",
			token:    "TestUseCase_Set1",
		},
		{
			name: "success2",
			args: args{
				ctx:   setToken(context.Background(), "TestUseCase_Set2"),
				key:   "343",
				value: "343",
			},
			token: "TestUseCase_Set2",
		},
		{
			name: "alreadyExists",
			args: args{
				ctx:   setToken(context.Background(), "TestUseCase_Set3"),
				key:   "343",
				value: "343",
			},
			token: "TestUseCase_Set3",
		},
		{
			name: "noToken",
			args: args{
				ctx:   context.Background(),
				key:   "343",
				value: "343",
			},
			token:   "TestUseCase_Set3",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				_, err := u.Register(tt.args.ctx, tt.username, tt.password)
				if err != nil && !errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				err = u.storage.AddToken(tt.args.ctx, tt.token, tt.username)
				if err != nil {
					t.Error(err)
				}
			}

			if err = u.Set(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, err := u.Get(tt.args.ctx, tt.args.key)
				if err != nil {
					t.Error(err)
				}

				if got != tt.args.value {
					t.Errorf("Set() got = %v, want %v", got, tt.args.value)
				}
			}
		})
	}
}

func TestUseCase_getUsernameFromContext(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		args         args
		token        string
		wantUsername string
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_getUsernameFromContext1"),
			},
			token:        "TestUseCase_getUsernameFromContext1",
			wantUsername: "XXX",
		},
		{
			name: "success2",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_getUsernameFromContext2"),
			},
			token:        "TestUseCase_getUsernameFromContext2",
			wantUsername: "wefwmfkmwmf",
		},
		{
			name: "noToken",
			args: args{
				ctx: context.Background(),
			},
			token:        "",
			wantUsername: "efgwegw",
			wantErr:      true,
		},
		{
			name: "errorNotFound",
			args: args{
				ctx: setToken(context.Background(), "TestUseCase_getUsernameFromContext3"),
			},
			token:        "TestUseCase_getUsernameFromContext3",
			wantUsername: "qwdqfdqfqf",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				err := u.storage.AddToken(tt.args.ctx, tt.token, tt.wantUsername)
				if err != nil {
					t.Error(err)
				}
			}

			gotUsername, err := u.getUsernameFromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUsernameFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if gotUsername != tt.wantUsername {
				t.Errorf("getUsernameFromContext() gotUsername = %v, want %v", gotUsername, tt.wantUsername)
			}
		})
	}
}

func TestUseCase_storeToken(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		token    string
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      setHeader(setToken(context.Background(), "TestUseCase_storeToken1")),
				token:    "TestUseCase_storeToken1",
				username: "XXX",
			},
		},
		{
			name: "success2",
			args: args{
				ctx:      setHeader(setToken(context.Background(), "TestUseCase_storeToken2")),
				token:    "TestUseCase_storeToken2",
				username: "XXwwegfwggwefX",
			},
		},
		{
			name: "noHeader",
			args: args{
				ctx:      context.Background(),
				token:    "",
				username: "XXwwegfwggwefX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}

			if !tt.wantErr {
				err = u.storage.AddToken(tt.args.ctx, tt.args.token, tt.args.username)
				if err != nil {
					t.Error(err)
				}
			}

			if err = u.storeToken(tt.args.ctx, tt.args.token, tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("storeToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, err := u.storage.GetUsername(tt.args.ctx, tt.args.token)
				if err != nil {
					t.Error(err)
				}
				if got != tt.args.username {
					t.Errorf("storeToken() got = %v, want %v", got, tt.args.username)
				}
			}
		})
	}
}

func TestUseCase_validateToken(t *testing.T) {
	store, err := storage.New(storage.Config{URI: ":800"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:   setToken(context.Background(), "TestUseCase_validateToken1"),
				token: "TestUseCase_validateToken1",
			},
			wantOk: true,
		},
		{
			name: "success2",
			args: args{
				ctx:   setToken(context.Background(), "TestUseCase_validateToken2"),
				token: "TestUseCase_validateToken2",
			},
			wantOk: true,
		},
		{
			name: "notFound",
			args: args{
				ctx:   setHeader(setToken(context.Background(), "TestUseCase_validateToken2")),
				token: "notFound",
			},
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store}
			if !tt.wantErr && tt.wantOk {
				err := store.AddToken(tt.args.ctx, tt.args.token, "XXX")
				if err != nil {
					t.Error(err)
				}
			}
			gotOk, err := u.validateToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("validateToken() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_generateToken(t *testing.T) {
	mp := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		token, err := generateToken()
		if err != nil {
			t.Error(err)
		}
		if _, ok := mp[token]; ok {
			t.Error("duplicate token")
		}
		mp[token] = struct{}{}
	}
}

func Test_getOrCreateToken(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantOk    bool
		wantToken string
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: setToken(context.Background(), "Test_getOrCreateToken1"),
			},
			wantOk:    true,
			wantToken: "Test_getOrCreateToken1",
		},
		{
			name: "success2",
			args: args{
				ctx: setToken(context.Background(), "Test_getOrCreateToken2"),
			},
			wantOk:    true,
			wantToken: "Test_getOrCreateToken2",
		},
		{
			name: "noToken",
			args: args{
				ctx: context.Background(),
			},
			wantOk: false,
		},
		{
			name: "noTokenInMetadata",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{})),
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantOk {
				md := metadata.New(map[string]string{"token": tt.wantToken})
				tt.args.ctx = metadata.NewIncomingContext(tt.args.ctx, md)
			}

			gotOk, gotToken, err := getOrCreateToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOrCreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("getOrCreateToken() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if !tt.wantOk || err != nil {
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("getOrCreateToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
