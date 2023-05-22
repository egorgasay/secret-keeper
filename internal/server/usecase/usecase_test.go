package usecase

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"secret-keeper/internal/server/storage"
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
				ctx:      setHeader(context.Background()),
				username: "XXX-16",
				password: "XXX-16",
			},
			token: "XXX-16",
		},

		{
			name: "success2",
			args: args{
				ctx:      setHeader(context.Background()),
				username: "admin3",
				password: "tgfqdf",
			},
			token: "qwerty",
		},

		{
			name: "error",
			args: args{
				ctx:      setHeader(context.Background()),
				username: "XX2424rt52dXXXX",
				password: "XwdefwegfXXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.token, nil
			}}

			var token string
			if !tt.wantErr {
				token, err = u.Register(tt.args.ctx, tt.args.username, tt.args.password)
				if err != nil && !errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				u.getOrCreateToken = func(ctx context.Context) (bool, string, error) {
					return true, tt.token, nil
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
				ctx: setHeader(context.Background()),
				key: "XXX-16",
			},
			username: "qwdq",
			password: "XXXX",
			token:    "XXX-16",
		},
		{
			name: "err",
			args: args{
				ctx: setHeader(context.Background()),
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
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.token, nil
			}}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
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
				ctx: setHeader(context.Background()),
				key: "XXX-16",
			},
			username: "qwdq",
			password: "XXXX",
			token:    "XXX-16",
			want:     "qwdq",
		},
		{
			name: "err",
			args: args{
				ctx: setHeader(context.Background()),
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
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.token, nil
			}}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
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
				ctx: setHeader(context.Background()),
			},
			username: "XXXX455",
			password: "XXXX456",
			token:    "XXX44-16",
			want: []string{
				"333", "444", "XXX44",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				ctx: setHeader(context.Background()),
			},
			username: "XXXX4355",
			password: "XXXX456",
			token:    "X32XX44-16",
			want:     []string{},
		},
		{
			name: "error",
			args: args{
				ctx: setHeader(context.Background()),
			},
			username: "X244355",
			password: "XXXX456",
			token:    "X32XX42544-16",
			want:     []string{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return false, tt.token, nil
			}}

			if !tt.wantErr {
				if _, err = u.Register(tt.args.ctx, tt.username, tt.password); err != nil &&
					!errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}

				u.getOrCreateToken = func(ctx context.Context) (bool, string, error) {
					return true, tt.token, nil
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
			if !isTheSameArray(got, tt.want) {
				t.Errorf("GetAllNames() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func isTheSameArray(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	var tmp = make(map[string]string, len(a))
	for _, el := range a {
		tmp[el] = el
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
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
				ctx:      setHeader(context.Background()),
				username: "XXX-16",
				password: "XXX-16",
			},
			token: "XXX-16",
		},

		{
			name: "success2",
			args: args{
				ctx:      setHeader(context.Background()),
				username: "admin3",
				password: "tgfqdf",
			},
			token: "qwerty",
		},

		{
			name: "error",
			args: args{
				ctx: setHeader(context.Background()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return false, tt.token, nil
			}}

			got, err := u.Register(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				if !errors.Is(err, storage.ErrAlreadyExists) {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !tt.wantErr {
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
				ctx:   setHeader(context.Background()),
				key:   "333",
				value: "333",
			},
			username: "qwe",
			password: "XXX",
			token:    "3334",
		},
		{
			name: "success2",
			args: args{
				ctx:   setHeader(context.Background()),
				key:   "343",
				value: "343",
			},
			token: "333445",
		},
		{
			name: "alreadyExists",
			args: args{
				ctx:   setHeader(context.Background()),
				key:   "343",
				value: "343",
			},
			token: "333445",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.token, nil
			}}

			if !tt.wantErr {
				_, err := u.Register(tt.args.ctx, tt.username, tt.password)
				if err != nil && !errors.Is(err, storage.ErrAlreadyExists) {
					t.Error(err)
				}
			}

			if err := u.Set(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
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
				ctx: setHeader(context.Background()),
			},
			token:        "qwe",
			wantUsername: "XXX",
		},
		{
			name: "success2",
			args: args{
				ctx: setHeader(context.Background()),
			},
			token:        "gmkwetm234j2n3fjn2j3nr",
			wantUsername: "wefwmfkmwmf",
		},
		{
			name: "noHeader",
			args: args{
				ctx: context.Background(),
			},
			token:        "fwe",
			wantUsername: "efgwegw",
			wantErr:      true,
		},
		{
			name: "errorNotFound",
			args: args{
				ctx: context.Background(),
			},
			token:        "qfwgqfe",
			wantUsername: "qwdqfdqfqf",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.token, nil
			}}

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
				ctx:      setHeader(context.Background()),
				token:    "qwe",
				username: "XXX",
			},
		},
		{
			name: "success2",
			args: args{
				ctx:      setHeader(context.Background()),
				token:    "qwdqwfqdwwefwge",
				username: "XXwwegfwggwefX",
			},
		},
		{
			name: "noHeader",
			args: args{
				ctx:      context.Background(),
				token:    "qwdqwfqdwwefwge",
				username: "XXwwegfwggwefX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCase{storage: store, getOrCreateToken: func(ctx context.Context) (bool, string, error) {
				return true, tt.args.token, nil
			}}
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
	type fields struct {
		storage *storage.Storage
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				storage: tt.fields.storage,
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
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateToken() got = %v, want %v", got, tt.want)
			}
		})
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotToken, err := getOrCreateToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOrCreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("getOrCreateToken() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotToken != tt.wantToken {
				t.Errorf("getOrCreateToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
