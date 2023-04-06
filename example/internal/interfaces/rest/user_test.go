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

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userInterceptor interceptors.UserInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want *UserHandler
	}{
		{
			name: "ok",
			args: args{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			want: &UserHandler{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.userInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewUserCreate(t)
	createjson, _ := json.Marshal(create)
	user := mock_models.NewUser(t)
	userjson, _ := json.Marshal(user)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
		wantBody   *bytes.Buffer
		wantStatus int
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(user, nil)
			},
			fields: fields{
				userInterceptor: userInterceptor,
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
			wantBody:   bytes.NewBuffer(userjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				userInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				userInterceptor: userInterceptor,
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
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.Create(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	arch := mock_models.NewArch(t)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
				userInterceptor.EXPECT().Delete(gomock.Any(), arch.ID, models.Guest).Return(nil)
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   &bytes.Buffer{},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "permission denied",
			setup: func() {
				userInterceptor.EXPECT().
					Delete(gomock.Any(), arch.ID, models.Guest).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
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
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(arch.ID))
			h.Delete(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	user := mock_models.NewUser(t)
	userjson, _ := json.Marshal(user)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
				userInterceptor.EXPECT().Get(gomock.Any(), user.ID, models.Guest).Return(user, nil)
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(userjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				userInterceptor.EXPECT().
					Get(gomock.Any(), user.ID, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
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
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(user.ID))
			h.Get(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &models.UserFilter{}
	listUsers := []*models.User{mock_models.NewUser(t)}
	listUsersJson, _ := json.Marshal(listUsers)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
				userInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(listUsers, uint64(len(listUsers)), nil)
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(listUsersJson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				userInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
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
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.List(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewUserUpdate(t)
	updatejson, _ := json.Marshal(update)
	user := mock_models.NewUser(t)
	userJson, _ := json.Marshal(user)
	type fields struct {
		userInterceptor interceptors.UserInterceptor
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
				userInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(user, nil)
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(userJson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				userInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &UserHandler{
				userInterceptor: tt.fields.userInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(update.ID))
			h.Update(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Create() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Create() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}