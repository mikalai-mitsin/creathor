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

func TestEquipmentHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	create := mock_models.NewEquipmentCreate(t)
	createjson, _ := json.Marshal(create)
	equipment := mock_models.NewEquipment(t)
	equipmentjson, _ := json.Marshal(equipment)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(equipment, nil)
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(createjson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(equipmentjson),
			wantStatus: http.StatusCreated,
		},
		{
			name: "permission denied",
			setup: func() {
				equipmentInterceptor.EXPECT().
					Create(gomock.Any(), create, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
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

func TestEquipmentHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor.EXPECT().
					Delete(gomock.Any(), equipment.ID, models.Guest).
					Return(nil)
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
				equipmentInterceptor.EXPECT().
					Delete(gomock.Any(), equipment.ID, models.Guest).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(equipment.ID))
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

func TestEquipmentHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	equipment := mock_models.NewEquipment(t)
	equipmentjson, _ := json.Marshal(equipment)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor.EXPECT().
					Get(gomock.Any(), equipment.ID, models.Guest).
					Return(equipment, nil)
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(equipmentjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				equipmentInterceptor.EXPECT().
					Get(gomock.Any(), equipment.ID, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.args.request
			ctx.AddParam("id", string(equipment.ID))
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

func TestEquipmentHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	filter := &models.EquipmentFilter{}
	listEquipment := []*models.Equipment{mock_models.NewEquipment(t)}
	listEquipmentjson, _ := json.Marshal(listEquipment)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(listEquipment, uint64(len(listEquipment)), nil)
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			args: args{
				request: (&http.Request{}).WithContext(
					context.WithValue(context.Background(), UserContextKey, models.Guest),
				),
			},
			wantBody:   bytes.NewBuffer(listEquipmentjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				equipmentInterceptor.EXPECT().
					List(gomock.Any(), filter, models.Guest).
					Return(nil, uint64(0), errs.NewPermissionDenied())
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
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

func TestEquipmentHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			args: args{
				router: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

func TestEquipmentHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	update := mock_models.NewEquipmentUpdate(t)
	updatejson, _ := json.Marshal(update)
	equipment := mock_models.NewEquipment(t)
	equipmentjson, _ := json.Marshal(equipment)
	type fields struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
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
				equipmentInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(equipment, nil)
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			args: args{
				request: (&http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(bytes.NewBuffer(updatejson)),
				}).WithContext(context.WithValue(context.Background(), UserContextKey, models.Guest)),
			},
			wantBody:   bytes.NewBuffer(equipmentjson),
			wantStatus: http.StatusOK,
		},
		{
			name: "permission denied",
			setup: func() {
				equipmentInterceptor.EXPECT().
					Update(gomock.Any(), update, models.Guest).
					Return(nil, errs.NewPermissionDenied())
			},
			fields: fields{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
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
			h := &EquipmentHandler{
				equipmentInterceptor: tt.fields.equipmentInterceptor,
				logger:               tt.fields.logger,
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

func TestNewEquipmentHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
	}
	tests := []struct {
		name string
		args args
		want *EquipmentHandler
	}{
		{
			name: "ok",
			args: args{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			want: &EquipmentHandler{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEquipmentHandler(tt.args.equipmentInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewEquipmentHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
