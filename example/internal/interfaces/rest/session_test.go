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

func TestSessionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewSessionCreate(t)
	createjson, _ := json.Marshal(create)
	session := mock_models.NewSession(t)
	sessionjson, _ := json.Marshal(session)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(session, nil)
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(sessionjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				sessionInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
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

func TestSessionHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	session := mock_models.NewSession(t)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor.EXPECT().
					Delete(gomock.Any(), session.ID, models.Guest).
					Return(nil)
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
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
				sessionInterceptor.EXPECT().
					Delete(gomock.Any(), session.ID, models.Guest).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
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
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(session.ID))
			h.Delete(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Delete() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Delete() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestSessionHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	session := mock_models.NewSession(t)
	sessionjson, _ := json.Marshal(session)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor.EXPECT().
					Get(gomock.Any(), session.ID, models.Guest).
					Return(session, nil)
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(sessionjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				sessionInterceptor.EXPECT().
					Get(gomock.Any(), session.ID, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
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
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(session.ID))
			h.Get(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Get() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Get() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestSessionHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &models.SessionFilter{}
	listSessions := []*models.Session{mock_models.NewSession(t)}
	listSessionsjson, _ := json.Marshal(listSessions)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(listSessions, uint64(len(listSessions)), nil)
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(listSessionsjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				sessionInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
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
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			h.List(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("List() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("List() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestSessionHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestSessionHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewSessionUpdate(t)
	updatejson, _ := json.Marshal(update)
	session := mock_models.NewSession(t)
	sessionjson, _ := json.Marshal(session)
	type fields struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
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
				sessionInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(session, nil)
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(sessionjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				sessionInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
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
			h := &SessionHandler{
				sessionInterceptor: tt.fields.sessionInterceptor,
				logger:             tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(update.ID))
			h.Update(ctx)
			if !reflect.DeepEqual(w.Code, tt.wantStatus) {
				t.Errorf("Update() gotStatus = %v, wantStatus %v", w.Code, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(w.Body, tt.wantBody) {
				t.Errorf("Update() gotBody = %v, wantBody %v", w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestNewSessionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
	}
	tests := []struct {
		name string
		args args
		want *SessionHandler
	}{
		{
			name: "ok",
			args: args{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			want: &SessionHandler{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionHandler(tt.args.sessionInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSessionHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
