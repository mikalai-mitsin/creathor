package interceptors

import (
	"context"
	"errors"
	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	mock_usecases "github.com/018bf/example/internal/domain/usecases/mock"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"syreclabs.com/go/faker"
	"testing"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

func TestNewEquipmentInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase      usecases.AuthUseCase
		equipmentUseCase usecases.EquipmentUseCase
		logger           log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.EquipmentInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			want: &EquipmentInterceptor{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewEquipmentInterceptor(tt.args.equipmentUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEquipmentInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		authUseCase      usecases.AuthUseCase
		equipmentUseCase usecases.EquipmentUseCase
		logger           log.Logger
	}
	type args struct {
		ctx         context.Context
		id          string
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDetail, equipment).
					Return(nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDetail, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Equipment not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	create := mock_models.NewEquipmentCreate(t)
	type fields struct {
		equipmentUseCase usecases.EquipmentUseCase
		authUseCase      usecases.AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.EquipmentCreate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentCreate, create).
					Return(nil)
				equipmentUseCase.EXPECT().Create(ctx, create).Return(equipment, nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentCreate, create).
					Return(nil)
				equipmentUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	update := mock_models.NewEquipmentUpdate(t)
	type fields struct {
		equipmentUseCase usecases.EquipmentUseCase
		authUseCase      usecases.AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.EquipmentUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().Update(ctx, update).Return(equipment, nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentUseCase usecases.EquipmentUseCase
		authUseCase      usecases.AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx         context.Context
		id          string
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDelete, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Equipment not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDelete, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDelete, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				id:          equipment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewEquipmentFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	equipment := make([]*models.Equipment, 0, count)
	for i := uint64(0); i < count; i++ {
		equipment = append(equipment, mock_models.NewEquipment(t))
	}
	type fields struct {
		equipmentUseCase usecases.EquipmentUseCase
		authUseCase      usecases.AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.EquipmentFilter
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Equipment
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentList, filter).
					Return(nil)
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(equipment, count, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    equipment,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentList, filter).
					Return(nil)
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("EquipmentInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
