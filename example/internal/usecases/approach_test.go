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
	"syreclabs.com/go/faker"
)

func TestNewApproachUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		approachRepository repositories.ApproachRepository
		clock              clock.Clock
		logger             log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.ApproachUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			want: &ApproachUseCase{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewApproachUseCase(tt.args.approachRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApproachUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	type fields struct {
		approachRepository repositories.ApproachRepository
		logger             log.Logger
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
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				approachRepository.EXPECT().Get(ctx, approach.ID).Return(approach, nil)
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
			},
			args: args{
				ctx: ctx,
				id:  approach.ID,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "Approach not found",
			setup: func() {
				approachRepository.EXPECT().Get(ctx, approach.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
			},
			args: args{
				ctx: ctx,
				id:  approach.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ApproachUseCase{
				approachRepository: tt.fields.approachRepository,
				logger:             tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var approaches []*models.Approach
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		approaches = append(approaches, mock_models.NewApproach(t))
	}
	filter := mock_models.NewApproachFilter(t)
	type fields struct {
		approachRepository repositories.ApproachRepository
		logger             log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ApproachFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Approach
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				approachRepository.EXPECT().List(ctx, filter).Return(approaches, nil)
				approachRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    approaches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				approachRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
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
				approachRepository.EXPECT().List(ctx, filter).Return(approaches, nil)
				approachRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
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
			u := &ApproachUseCase{
				approachRepository: tt.fields.approachRepository,
				logger:             tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ApproachUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestApproachUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewApproachCreate(t)
	now := time.Now().UTC()
	type fields struct {
		approachRepository repositories.ApproachRepository
		clock              clock.Clock
		logger             log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.ApproachCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				approachRepository.EXPECT().
					Create(
						ctx,
						&models.Approach{
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Approach{
				ID:        "",
				UpdatedAt: now,
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				approachRepository.EXPECT().
					Create(
						ctx,
						&models.Approach{
							ID:        "",
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		// TODO: Add validation rules or delete this case
		//{
		//	name: "invalid",
		//	setup: func() {
		//	},
		//	fields: fields{
		//		approachRepository: approachRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.ApproachCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ApproachUseCase{
				approachRepository: tt.fields.approachRepository,
				clock:              tt.fields.clock,
				logger:             tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewApproachUpdate(t)
	now := approach.UpdatedAt
	type fields struct {
		approachRepository repositories.ApproachRepository
		clock              clock.Clock
		logger             log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.ApproachUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				approachRepository.EXPECT().
					Get(ctx, update.ID).Return(approach, nil)
				approachRepository.EXPECT().
					Update(ctx, approach).Return(nil)
			},
			fields: fields{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				approachRepository.EXPECT().
					Get(ctx, update.ID).
					Return(approach, nil)
				approachRepository.EXPECT().
					Update(ctx, approach).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Approach not found",
			setup: func() {
				approachRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
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
				approachRepository: approachRepository,
				clock:              clockMock,
				logger:             logger,
			},
			args: args{
				ctx: ctx,
				update: &models.ApproachUpdate{
					ID: faker.Number().Number(1),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ApproachUseCase{
				approachRepository: tt.fields.approachRepository,
				clock:              tt.fields.clock,
				logger:             tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	approachRepository := mock_repositories.NewMockApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	type fields struct {
		approachRepository repositories.ApproachRepository
		logger             log.Logger
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				approachRepository.EXPECT().
					Delete(ctx, approach.ID).
					Return(nil)
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
			},
			args: args{
				ctx: ctx,
				id:  approach.ID,
			},
			wantErr: nil,
		},
		{
			name: "Approach not found",
			setup: func() {
				approachRepository.EXPECT().
					Delete(ctx, approach.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				approachRepository: approachRepository,
				logger:             logger,
			},
			args: args{
				ctx: ctx,
				id:  approach.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ApproachUseCase{
				approachRepository: tt.fields.approachRepository,
				logger:             tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
