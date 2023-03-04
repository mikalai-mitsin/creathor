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

func TestNewDayUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		dayRepository repositories.DayRepository
		clock         clock.Clock
		logger        log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.DayUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
			want: &DayUseCase{
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewDayUseCase(tt.args.dayRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		dayRepository repositories.DayRepository
		logger        log.Logger
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
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayRepository.EXPECT().Get(ctx, day.ID).Return(day, nil)
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "Day not found",
			setup: func() {
				dayRepository.EXPECT().Get(ctx, day.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &DayUseCase{
				dayRepository: tt.fields.dayRepository,
				logger:        tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listDays []*models.Day
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listDays = append(listDays, mock_models.NewDay(t))
	}
	filter := mock_models.NewDayFilter(t)
	type fields struct {
		dayRepository repositories.DayRepository
		logger        log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.DayFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Day
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayRepository.EXPECT().List(ctx, filter).Return(listDays, nil)
				dayRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listDays,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				dayRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
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
				dayRepository.EXPECT().List(ctx, filter).Return(listDays, nil)
				dayRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
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
			u := &DayUseCase{
				dayRepository: tt.fields.dayRepository,
				logger:        tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DayUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDayUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewDayCreate(t)
	now := time.Now().UTC()
	type fields struct {
		dayRepository repositories.DayRepository
		clock         clock.Clock
		logger        log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.DayCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				dayRepository.EXPECT().
					Create(
						ctx,
						&models.Day{
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
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Day{
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
				dayRepository.EXPECT().
					Create(
						ctx,
						&models.Day{
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
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
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
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx:    ctx,
				create: &models.DayCreate{},
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
			u := &DayUseCase{
				dayRepository: tt.fields.dayRepository,
				clock:         tt.fields.clock,
				logger:        tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewDayUpdate(t)
	now := day.UpdatedAt
	type fields struct {
		dayRepository repositories.DayRepository
		clock         clock.Clock
		logger        log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.DayUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				dayRepository.EXPECT().
					Get(ctx, update.ID).Return(day, nil)
				dayRepository.EXPECT().
					Update(ctx, day).Return(nil)
			},
			fields: fields{
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				dayRepository.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				dayRepository.EXPECT().
					Update(ctx, day).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Day not found",
			setup: func() {
				dayRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
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
				dayRepository: dayRepository,
				clock:         clockMock,
				logger:        logger,
			},
			args: args{
				ctx: ctx,
				update: &models.DayUpdate{
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
			u := &DayUseCase{
				dayRepository: tt.fields.dayRepository,
				clock:         tt.fields.clock,
				logger:        tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayRepository := mock_repositories.NewMockDayRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		dayRepository repositories.DayRepository
		logger        log.Logger
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
				dayRepository.EXPECT().
					Delete(ctx, day.ID).
					Return(nil)
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: nil,
		},
		{
			name: "Day not found",
			setup: func() {
				dayRepository.EXPECT().
					Delete(ctx, day.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				dayRepository: dayRepository,
				logger:        logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &DayUseCase{
				dayRepository: tt.fields.dayRepository,
				logger:        tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("DayUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
