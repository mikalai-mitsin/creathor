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

func TestNewUserSessionUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userSessionRepository repositories.UserSessionRepository
		clock                 clock.Clock
		logger                log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.UserSessionUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
			want: &UserSessionUseCase{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserSessionUseCase(tt.args.userSessionRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserSessionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	type fields struct {
		userSessionRepository repositories.UserSessionRepository
		logger                log.Logger
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
		want    *models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userSessionRepository.EXPECT().Get(ctx, userSession.ID).Return(userSession, nil)
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
			},
			args: args{
				ctx: ctx,
				id:  userSession.ID,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "UserSession not found",
			setup: func() {
				userSessionRepository.EXPECT().Get(ctx, userSession.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
			},
			args: args{
				ctx: ctx,
				id:  userSession.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserSessionUseCase{
				userSessionRepository: tt.fields.userSessionRepository,
				logger:                tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var userSessions []*models.UserSession
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		userSessions = append(userSessions, mock_models.NewUserSession(t))
	}
	filter := mock_models.NewUserSessionFilter(t)
	type fields struct {
		userSessionRepository repositories.UserSessionRepository
		logger                log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserSessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.UserSession
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userSessionRepository.EXPECT().List(ctx, filter).Return(userSessions, nil)
				userSessionRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    userSessions,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				userSessionRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
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
				userSessionRepository.EXPECT().List(ctx, filter).Return(userSessions, nil)
				userSessionRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
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
			u := &UserSessionUseCase{
				userSessionRepository: tt.fields.userSessionRepository,
				logger:                tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserSessionUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserSessionUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewUserSessionCreate(t)
	now := time.Now().UTC()
	type fields struct {
		userSessionRepository repositories.UserSessionRepository
		clock                 clock.Clock
		logger                log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.UserSessionCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userSessionRepository.EXPECT().
					Create(
						ctx,
						&models.UserSession{
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.UserSession{
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
				userSessionRepository.EXPECT().
					Create(
						ctx,
						&models.UserSession{
							ID:        "",
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
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
		//		userSessionRepository: userSessionRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.UserSessionCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserSessionUseCase{
				userSessionRepository: tt.fields.userSessionRepository,
				clock:                 tt.fields.clock,
				logger:                tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewUserSessionUpdate(t)
	now := userSession.UpdatedAt
	type fields struct {
		userSessionRepository repositories.UserSessionRepository
		clock                 clock.Clock
		logger                log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.UserSessionUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userSessionRepository.EXPECT().
					Get(ctx, update.ID).Return(userSession, nil)
				userSessionRepository.EXPECT().
					Update(ctx, userSession).Return(nil)
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userSessionRepository.EXPECT().
					Get(ctx, update.ID).
					Return(userSession, nil)
				userSessionRepository.EXPECT().
					Update(ctx, userSession).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "UserSession not found",
			setup: func() {
				userSessionRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
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
				userSessionRepository: userSessionRepository,
				clock:                 clockMock,
				logger:                logger,
			},
			args: args{
				ctx: ctx,
				update: &models.UserSessionUpdate{
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
			u := &UserSessionUseCase{
				userSessionRepository: tt.fields.userSessionRepository,
				clock:                 tt.fields.clock,
				logger:                tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSessionRepository := mock_repositories.NewMockUserSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	type fields struct {
		userSessionRepository repositories.UserSessionRepository
		logger                log.Logger
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
				userSessionRepository.EXPECT().
					Delete(ctx, userSession.ID).
					Return(nil)
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
			},
			args: args{
				ctx: ctx,
				id:  userSession.ID,
			},
			wantErr: nil,
		},
		{
			name: "UserSession not found",
			setup: func() {
				userSessionRepository.EXPECT().
					Delete(ctx, userSession.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				userSessionRepository: userSessionRepository,
				logger:                logger,
			},
			args: args{
				ctx: ctx,
				id:  userSession.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserSessionUseCase{
				userSessionRepository: tt.fields.userSessionRepository,
				logger:                tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
