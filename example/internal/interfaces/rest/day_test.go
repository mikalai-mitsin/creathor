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

func TestDayHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewDayCreate(t)
	createjson, _ := json.Marshal(create)
	day := mock_models.NewDay(t)
	dayjson, _ := json.Marshal(day)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor.EXPECT().Create(gomock.Any(), create, models.Guest).Return(day, nil)
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(dayjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				dayInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
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

func TestDayHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	day := mock_models.NewDay(t)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor.EXPECT().Delete(gomock.Any(), day.ID, models.Guest).Return(nil)
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
				dayInterceptor.EXPECT().
					Delete(gomock.Any(), day.ID, models.Guest).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(day.ID))
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

func TestDayHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	day := mock_models.NewDay(t)
	dayjson, _ := json.Marshal(day)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor.EXPECT().Get(gomock.Any(), day.ID, models.Guest).Return(day, nil)
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(dayjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				dayInterceptor.EXPECT().
					Get(gomock.Any(), day.ID, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(day.ID))
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

func TestDayHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &models.DayFilter{}
	listDays := []*models.Day{mock_models.NewDay(t)}
	listDaysjson, _ := json.Marshal(listDays)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(listDays, uint64(len(listDays)), nil)
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(listDaysjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				dayInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
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

func TestDayHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestDayHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewDayUpdate(t)
	updatejson, _ := json.Marshal(update)
	day := mock_models.NewDay(t)
	dayjson, _ := json.Marshal(day)
	type fields struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
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
				dayInterceptor.EXPECT().Update(gomock.Any(), update, models.Guest).Return(day, nil)
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(dayjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				dayInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				dayInterceptor: dayInterceptor,
				logger:         logger,
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
			h := &DayHandler{
				dayInterceptor: tt.fields.dayInterceptor,
				logger:         tt.fields.logger,
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

func TestNewDayHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_interceptors.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		dayInterceptor interceptors.DayInterceptor
		logger         log.Logger
	}
	tests := []struct {
		name string
		args args
		want *DayHandler
	}{
		{
			name: "ok",
			args: args{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			want: &DayHandler{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDayHandler(tt.args.dayInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
