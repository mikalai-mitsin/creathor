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

func TestNewUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userRepository repositories.UserRepository
		clock          clock.Clock
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.UserUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				userRepository: userRepository,
				clock:          clockMock,
				logger:         logger,
			},
			want: &UserUseCase{
				userRepository: userRepository,
				clock:          clockMock,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserUseCase(tt.args.userRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userRepository repositories.UserRepository
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
		want    *models.User
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				userRepository.EXPECT().Get(ctx, user.ID).Return(user, nil)
			},
			fields: fields{
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				userRepository.EXPECT().Get(ctx, user.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserUseCase{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var users []*models.User
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		users = append(users, mock_models.NewUser(t))
	}
	filter := mock_models.NewUserFilter(t)
	type fields struct {
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.User
		want1   uint64
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				userRepository.EXPECT().List(ctx, filter).Return(users, nil)
				userRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    users,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				userRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: userRepository,
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
				userRepository.EXPECT().List(ctx, filter).Return(users, nil)
				userRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: userRepository,
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
			u := &UserUseCase{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewUserCreate(t)
	now := time.Now().UTC()
	type fields struct {
		userRepository repositories.UserRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.UserCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userRepository.EXPECT().
					Create(
						ctx,
						&models.User{
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				userRepository: userRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.User{
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
				userRepository.EXPECT().
					Create(
						ctx,
						&models.User{
							ID:        "",
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: userRepository,
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
		// {
		//     name: "invalid",
		//     setup: func() {
		//     },
		//     fields: fields{
		//         userRepository: userRepository,
		//         logger:           logger,
		//     },
		//     args: args{
		//         ctx: ctx,
		//         create: &models.UserCreate{},
		//     },
		//     want: nil,
		//     wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserUseCase{
				userRepository: tt.fields.userRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewUserUpdate(t)
	now := user.UpdatedAt
	type fields struct {
		userRepository repositories.UserRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.UserUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userRepository.EXPECT().
					Get(ctx, update.ID).Return(user, nil)
				userRepository.EXPECT().
					Update(ctx, user).Return(nil)
			},
			fields: fields{
				userRepository: userRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				userRepository.EXPECT().
					Get(ctx, update.ID).
					Return(user, nil)
				userRepository.EXPECT().
					Update(ctx, user).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: userRepository,
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
			name: "User not found",
			setup: func() {
				userRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userRepository: userRepository,
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
				userRepository: userRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				update: &models.UserUpdate{
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
			u := &UserUseCase{
				userRepository: tt.fields.userRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userRepository repositories.UserRepository
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
		wantErr *errs.Error
	}{
		{
			name: "ok",
			setup: func() {
				userRepository.EXPECT().
					Delete(ctx, user.ID).
					Return(nil)
			},
			fields: fields{
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				userRepository.EXPECT().
					Delete(ctx, user.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserUseCase{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
