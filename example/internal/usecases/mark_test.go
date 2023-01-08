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

func TestNewMarkUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		markRepository repositories.MarkRepository
		clock          clock.Clock
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.MarkUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				markRepository: markRepository,
				clock:          clockMock,
				logger:         logger,
			},
			want: &MarkUseCase{
				markRepository: markRepository,
				clock:          clockMock,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewMarkUseCase(tt.args.markRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarkUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	mark := mock_models.NewMark(t)
	type fields struct {
		markRepository repositories.MarkRepository
		logger         log.Logger
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
		want    *models.Mark
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				markRepository.EXPECT().Get(ctx, mark.ID).Return(mark, nil)
			},
			fields: fields{
				markRepository: markRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  mark.ID,
			},
			want:    mark,
			wantErr: nil,
		},
		{
			name: "Mark not found",
			setup: func() {
				markRepository.EXPECT().Get(ctx, mark.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				markRepository: markRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  mark.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &MarkUseCase{
				markRepository: tt.fields.markRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var marks []*models.Mark
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		marks = append(marks, mock_models.NewMark(t))
	}
	filter := mock_models.NewMarkFilter(t)
	type fields struct {
		markRepository repositories.MarkRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.MarkFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Mark
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				markRepository.EXPECT().List(ctx, filter).Return(marks, nil)
				markRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				markRepository: markRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    marks,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				markRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				markRepository: markRepository,
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
				markRepository.EXPECT().List(ctx, filter).Return(marks, nil)
				markRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				markRepository: markRepository,
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
			u := &MarkUseCase{
				markRepository: tt.fields.markRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MarkUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMarkUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewMarkCreate(t)
	now := time.Now().UTC()
	type fields struct {
		markRepository repositories.MarkRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.MarkCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Mark
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				markRepository.EXPECT().
					Create(
						ctx,
						&models.Mark{
							Name:      create.Name,
							Title:     create.Title,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				markRepository: markRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Mark{
				ID:        "",
				Name:      create.Name,
				Title:     create.Title,
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
				markRepository.EXPECT().
					Create(
						ctx,
						&models.Mark{
							ID:        "",
							Name:      create.Name,
							Title:     create.Title,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				markRepository: markRepository,
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
		// TODO: Add validation rules or delete this case
		//{
		//	name: "invalid",
		//	setup: func() {
		//	},
		//	fields: fields{
		//		markRepository: markRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.MarkCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &MarkUseCase{
				markRepository: tt.fields.markRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	mark := mock_models.NewMark(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewMarkUpdate(t)
	now := mark.UpdatedAt
	type fields struct {
		markRepository repositories.MarkRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.MarkUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Mark
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				markRepository.EXPECT().
					Get(ctx, update.ID).Return(mark, nil)
				markRepository.EXPECT().
					Update(ctx, mark).Return(nil)
			},
			fields: fields{
				markRepository: markRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    mark,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				markRepository.EXPECT().
					Get(ctx, update.ID).
					Return(mark, nil)
				markRepository.EXPECT().
					Update(ctx, mark).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				markRepository: markRepository,
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
			name: "Mark not found",
			setup: func() {
				markRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				markRepository: markRepository,
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
				markRepository: markRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				update: &models.MarkUpdate{
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
			u := &MarkUseCase{
				markRepository: tt.fields.markRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	markRepository := mock_repositories.NewMockMarkRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	mark := mock_models.NewMark(t)
	type fields struct {
		markRepository repositories.MarkRepository
		logger         log.Logger
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
				markRepository.EXPECT().
					Delete(ctx, mark.ID).
					Return(nil)
			},
			fields: fields{
				markRepository: markRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  mark.ID,
			},
			wantErr: nil,
		},
		{
			name: "Mark not found",
			setup: func() {
				markRepository.EXPECT().
					Delete(ctx, mark.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				markRepository: markRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  mark.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &MarkUseCase{
				markRepository: tt.fields.markRepository,
				logger:         tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
