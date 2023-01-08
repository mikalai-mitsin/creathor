package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	mock_interceptors "github.com/018bf/example/internal/domain/interceptors/mock"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/pkg/log"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"google.golang.org/grpc/metadata"
)

type Transport struct {
	MethodPath string
}

func (t Transport) Method() string {
	return t.MethodPath
}

func (t Transport) SetHeader(_ metadata.MD) error {
	return nil
}

func (t Transport) SendHeader(_ metadata.MD) error {
	return nil
}

func (t Transport) SetTrailer(_ metadata.MD) error {
	return nil
}

func TestAuthMiddleware_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	ctx := context.Background()
	token := mock_models.NewToken(t)
	user := mock_models.NewUser(t)
	ctxWithToken := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token.String()),
	}))
	ctxWithBadToken := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", "token.String()"),
	}))
	type fields struct {
		authInterceptor interceptors.AuthInterceptor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    context.Context
		wantErr error
		setup   func()
	}{
		{
			name: "ok",
			setup: func() {
				authInterceptor.EXPECT().Auth(ctxWithToken, token).Return(user, nil)
			},
			fields: fields{
				authInterceptor: authInterceptor,
			},
			args: args{
				ctx: ctxWithToken,
			},
			want:    context.WithValue(ctxWithToken, UserKey, user),
			wantErr: nil,
		},
		{
			name: "bad token",
			setup: func() {
				authInterceptor.EXPECT().
					Auth(ctxWithBadToken, models.Token("token.String()")).
					Return(nil, errs.NewBadToken()).
					Times(1)
			},
			fields: fields{
				authInterceptor: authInterceptor,
			},
			args: args{
				ctx: ctxWithBadToken,
			},
			want:    nil,
			wantErr: decodeError(errs.NewBadToken()),
		},
		{
			name:  "without token",
			setup: func() {},
			fields: fields{
				authInterceptor: authInterceptor,
			},
			args: args{
				ctx: ctx,
			},
			want:    context.WithValue(ctx, UserKey, models.Guest),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			m := AuthMiddleware{
				authInterceptor: tt.fields.authInterceptor,
			}
			got, err := m.Auth(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	type args struct {
		authInterceptor interceptors.AuthInterceptor
		logger          log.Logger
		config          *configs.Config
	}
	tests := []struct {
		name string
		args args
		want *AuthMiddleware
	}{
		{
			name: "ok",
			args: args{
				authInterceptor: authInterceptor,
			},
			want: &AuthMiddleware{
				authInterceptor: authInterceptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthMiddleware(tt.args.authInterceptor, tt.args.logger, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
