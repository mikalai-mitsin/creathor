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

func TestNewPlanUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		planRepository repositories.PlanRepository
		clock          clock.Clock
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.PlanUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
			want: &PlanUseCase{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPlanUseCase(tt.args.planRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	type fields struct {
		planRepository repositories.PlanRepository
		logger         log.Logger
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
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planRepository.EXPECT().Get(ctx, plan.ID).Return(plan, nil)
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  plan.ID,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "Plan not found",
			setup: func() {
				planRepository.EXPECT().Get(ctx, plan.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  plan.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PlanUseCase{
				planRepository: tt.fields.planRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listPlans []*models.Plan
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listPlans = append(listPlans, mock_models.NewPlan(t))
	}
	filter := mock_models.NewPlanFilter(t)
	type fields struct {
		planRepository repositories.PlanRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.PlanFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Plan
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planRepository.EXPECT().List(ctx, filter).Return(listPlans, nil)
				planRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listPlans,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				planRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
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
				planRepository.EXPECT().List(ctx, filter).Return(listPlans, nil)
				planRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
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
			u := &PlanUseCase{
				planRepository: tt.fields.planRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlanUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPlanUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewPlanCreate(t)
	now := time.Now().UTC()
	type fields struct {
		planRepository repositories.PlanRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.PlanCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				planRepository.EXPECT().
					Create(
						ctx,
						&models.Plan{
							Name:        create.Name,
							Repeat:      create.Repeat,
							EquipmentID: create.EquipmentID,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(nil)
			},
			fields: fields{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Plan{
				ID:          "",
				Name:        create.Name,
				Repeat:      create.Repeat,
				EquipmentID: create.EquipmentID,
				UpdatedAt:   now,
				CreatedAt:   now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				planRepository.EXPECT().
					Create(
						ctx,
						&models.Plan{
							ID:          "",
							Name:        create.Name,
							Repeat:      create.Repeat,
							EquipmentID: create.EquipmentID,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
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
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: &models.PlanCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(map[string]string{
				"name":         "cannot be blank",
				"repeat":       "cannot be blank",
				"equipment_id": "cannot be blank",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PlanUseCase{
				planRepository: tt.fields.planRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewPlanUpdate(t)
	now := plan.UpdatedAt
	type fields struct {
		planRepository repositories.PlanRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.PlanUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				planRepository.EXPECT().
					Get(ctx, update.ID).Return(plan, nil)
				planRepository.EXPECT().
					Update(ctx, plan).Return(nil)
			},
			fields: fields{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				planRepository.EXPECT().
					Get(ctx, update.ID).
					Return(plan, nil)
				planRepository.EXPECT().
					Update(ctx, plan).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Plan not found",
			setup: func() {
				planRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
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
				planRepository: planRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				update: &models.PlanUpdate{
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
			u := &PlanUseCase{
				planRepository: tt.fields.planRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planRepository := mock_repositories.NewMockPlanRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	type fields struct {
		planRepository repositories.PlanRepository
		logger         log.Logger
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
				planRepository.EXPECT().
					Delete(ctx, plan.ID).
					Return(nil)
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  plan.ID,
			},
			wantErr: nil,
		},
		{
			name: "Plan not found",
			setup: func() {
				planRepository.EXPECT().
					Delete(ctx, plan.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				planRepository: planRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  plan.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PlanUseCase{
				planRepository: tt.fields.planRepository,
				logger:         tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
