package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	mock_repositories "github.com/018bf/example/internal/domain/repositories/mock"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	mock_clock "github.com/018bf/example/pkg/clock/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
)

func TestNewEquipmentUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		equipmentRepository repositories.EquipmentRepository
		clock               clock.Clock
		logger              log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.EquipmentUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			want: &EquipmentUseCase{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewEquipmentUseCase(tt.args.equipmentRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewEquipmentUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentRepository repositories.EquipmentRepository
		logger              log.Logger
	}
	type args struct {
		ctx context.Context
		id  models.UUID
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
				equipmentRepository.EXPECT().Get(ctx, equipment.ID).Return(equipment, nil)
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
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
				equipmentRepository.EXPECT().
					Get(ctx, equipment.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
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
			u := &EquipmentUseCase{
				equipmentRepository: tt.fields.equipmentRepository,
				logger:              tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listEquipment []*models.Equipment
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listEquipment = append(listEquipment, mock_models.NewEquipment(t))
	}
	filter := mock_models.NewEquipmentFilter(t)
	type fields struct {
		equipmentRepository repositories.EquipmentRepository
		logger              log.Logger
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentRepository.EXPECT().List(ctx, filter).Return(listEquipment, nil)
				equipmentRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listEquipment,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				equipmentRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "count error",
			setup: func() {
				equipmentRepository.EXPECT().List(ctx, filter).Return(listEquipment, nil)
				equipmentRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EquipmentUseCase{
				equipmentRepository: tt.fields.equipmentRepository,
				logger:              tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("EquipmentUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestEquipmentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewEquipmentCreate(t)
	now := time.Now().UTC()
	type fields struct {
		equipmentRepository repositories.EquipmentRepository
		clock               clock.Clock
		logger              log.Logger
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				equipmentRepository.EXPECT().
					Create(
						ctx,
						&models.Equipment{
							Name:      create.Name,
							Repeat:    create.Repeat,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Equipment{
				ID:        "",
				Name:      create.Name,
				Repeat:    create.Repeat,
				Weight:    create.Weight,
				UpdatedAt: now,
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				equipmentRepository.EXPECT().
					Create(
						ctx,
						&models.Equipment{
							ID:        "",
							Name:      create.Name,
							Repeat:    create.Repeat,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				create: &models.EquipmentCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(map[string]string{
				"name":   "cannot be blank",
				"repeat": "cannot be blank",
				"weight": "cannot be blank",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EquipmentUseCase{
				equipmentRepository: tt.fields.equipmentRepository,
				clock:               tt.fields.clock,
				logger:              tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewEquipmentUpdate(t)
	now := equipment.UpdatedAt
	type fields struct {
		equipmentRepository repositories.EquipmentRepository
		clock               clock.Clock
		logger              log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.EquipmentUpdate
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
				clockMock.EXPECT().Now().Return(now)
				equipmentRepository.EXPECT().
					Get(ctx, update.ID).Return(equipment, nil)
				equipmentRepository.EXPECT().
					Update(ctx, equipment).Return(nil)
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
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
				clockMock.EXPECT().Now().Return(now)
				equipmentRepository.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				equipmentRepository.EXPECT().
					Update(ctx, equipment).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Equipment not found",
			setup: func() {
				equipmentRepository.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				clock:               clockMock,
				logger:              logger,
			},
			args: args{
				ctx: ctx,
				update: &models.EquipmentUpdate{
					ID: models.UUID("baduuid"),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EquipmentUseCase{
				equipmentRepository: tt.fields.equipmentRepository,
				clock:               tt.fields.clock,
				logger:              tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentRepository := mock_repositories.NewMockEquipmentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentRepository repositories.EquipmentRepository
		logger              log.Logger
	}
	type args struct {
		ctx context.Context
		id  models.UUID
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
				equipmentRepository.EXPECT().
					Delete(ctx, equipment.ID).
					Return(nil)
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: nil,
		},
		{
			name: "Equipment not found",
			setup: func() {
				equipmentRepository.EXPECT().
					Delete(ctx, equipment.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				equipmentRepository: equipmentRepository,
				logger:              logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &EquipmentUseCase{
				equipmentRepository: tt.fields.equipmentRepository,
				logger:              tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
