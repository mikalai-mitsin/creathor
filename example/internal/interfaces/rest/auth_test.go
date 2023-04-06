package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	mock_interceptors "github.com/018bf/example/internal/domain/interceptors/mock"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAuthHandler_CreateTokenPair(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewLogin(t)
	createjson, _ := json.Marshal(create)
	pair := mock_models.NewTokenPair(t)
	pairjson, _ := json.Marshal(pair)
	type fields struct {
		authInterceptor interceptors.AuthInterceptor
		logger          log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{
		{
			name: "ok",
			setup: func() {
				authInterceptor.EXPECT().
					CreateToken(gomock.Any(), create).
					Return(pair, nil)
			},
			fields: fields{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(pairjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				authInterceptor.EXPECT().
					CreateToken(gomock.Any(), create).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &AuthHandler{
				authInterceptor: tt.fields.authInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.CreateTokenPair(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("CreateTokenPair() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("CreateTokenPair() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestAuthHandler_RefreshTokenPair(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := &Refresh{
		Token: mock_models.NewToken(t),
	}
	createjson, _ := json.Marshal(create)
	pair := mock_models.NewTokenPair(t)
	pairjson, _ := json.Marshal(pair)
	type fields struct {
		authInterceptor interceptors.AuthInterceptor
		logger          log.Logger
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name       string
		setup      func()
		fields     fields
		args       args
		wantStatus int
		wantBody   *bytes.Buffer
	}{
		{
			name: "ok",
			setup: func() {
				authInterceptor.EXPECT().
					RefreshToken(gomock.Any(), create.Token).
					Return(pair, nil)
			},
			fields: fields{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(pairjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				authInterceptor.EXPECT().
					RefreshToken(gomock.Any(), create.Token).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &AuthHandler{
				authInterceptor: tt.fields.authInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.RefreshTokenPair(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("RefreshTokenPair() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("RefreshTokenPair() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		authInterceptor interceptors.AuthInterceptor
		logger          log.Logger
	}
	type args struct {
		router *gin.RouterGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				authInterceptor: tt.fields.authInterceptor,
				logger:          tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestNewAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_interceptors.NewMockAuthInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authInterceptor interceptors.AuthInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want *AuthHandler
	}{
		{
			name: "ok",
			args: args{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
			want: &AuthHandler{
				authInterceptor: authInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthHandler(tt.args.authInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}