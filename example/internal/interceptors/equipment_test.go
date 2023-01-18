package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	mock_usecases "github.com/018bf/example/internal/domain/usecases/mock"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

func TestNewEquipmentInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
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
				logger:           logger,
			},
			want: &EquipmentInterceptor{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewEquipmentInterceptor(tt.args.equipmentUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEquipmentInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentUseCase usecases.EquipmentUseCase
		logger           log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "Equipment not found",
			setup: func() {
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
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
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
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
		ctx    context.Context
		create *models.EquipmentCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentUseCase.EXPECT().Create(ctx, create).Return(equipment, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				equipmentUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
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
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
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
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentUseCase.EXPECT().Update(ctx, update).Return(equipment, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				equipmentUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
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
	equipmentUseCase := mock_usecases.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
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
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				logger:           tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
		logger           log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.EquipmentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Equipment
		want1   uint64
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(equipment, count, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    equipment,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
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
				logger:           tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
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
