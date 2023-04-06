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

func TestPlanHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewPlanCreate(t)
	createjson, _ := json.Marshal(create)
	plan := mock_models.NewPlan(t)
	planjson, _ := json.Marshal(plan)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(plan, nil)
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			wantBody:   bytes.NewBuffer(planjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				planInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			wantBody:   bytes.NewBufferString(errs.NewPermissionDenied().Error()),
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
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

func TestPlanHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	plan := mock_models.NewPlan(t)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor.EXPECT().Delete(gomock.Any(), plan.ID, models.Guest).Return(nil)
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
				planInterceptor.EXPECT().
					Delete(gomock.Any(), plan.ID, models.Guest).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(plan.ID))
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

func TestPlanHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	plan := mock_models.NewPlan(t)
	planjson, _ := json.Marshal(plan)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor.EXPECT().Get(gomock.Any(), plan.ID, models.Guest).Return(plan, nil)
			},
			fields: fields{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(planjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				planInterceptor.EXPECT().
					Get(gomock.Any(), plan.ID, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
				logger:          tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(plan.ID))
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

func TestPlanHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &models.PlanFilter{}
	listPlans := []*models.Plan{mock_models.NewPlan(t)}
	listPlansjson, _ := json.Marshal(listPlans)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(listPlans, uint64(len(listPlans)), nil)
			},
			fields: fields{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(listPlansjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				planInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
				logger:          tt.fields.logger,
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

func TestPlanHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
				logger:          tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestPlanHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewPlanUpdate(t)
	updatejson, _ := json.Marshal(update)
	plan := mock_models.NewPlan(t)
	planjson, _ := json.Marshal(plan)
	type fields struct {
		planInterceptor interceptors.PlanInterceptor
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
				planInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(plan, nil)
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			wantBody:   bytes.NewBuffer(planjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				planInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				planInterceptor: planInterceptor,
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
			h := &PlanHandler{
				planInterceptor: tt.fields.planInterceptor,
				logger:          tt.fields.logger,
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

func TestNewPlanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		planInterceptor interceptors.PlanInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want *PlanHandler
	}{
		{
			name: "ok",
			args: args{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			want: &PlanHandler{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlanHandler(tt.args.planInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
