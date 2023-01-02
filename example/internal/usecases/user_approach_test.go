package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/creathor/internal/domain/errs"
	"github.com/018bf/creathor/internal/domain/models"
	mock_models "github.com/018bf/creathor/internal/domain/models/mock"
	"github.com/018bf/creathor/internal/domain/repositories"
	mock_repositories "github.com/018bf/creathor/internal/domain/repositories/mock"
	"github.com/018bf/creathor/internal/domain/usecases"
	"github.com/018bf/creathor/pkg/log"
	mock_log "github.com/018bf/creathor/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"
)

func TestNewUserApproachUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.UserApproachUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			want: &UserApproachUseCase{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserApproachUseCase(tt.args.userApproachRepository, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserApproachUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	type fields struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
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
		want    *models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userApproachRepository.EXPECT().Get(ctx, userApproach.ID).Return(userApproach, nil)
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx: ctx,
				id:  userApproach.ID,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "UserApproach not found",
			setup: func() {
				userApproachRepository.EXPECT().Get(ctx, userApproach.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx: ctx,
				id:  userApproach.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserApproachUseCase{
				userApproachRepository: tt.fields.userApproachRepository,
				logger:                 tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var userApproaches []*models.UserApproach
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		userApproaches = append(userApproaches, mock_models.NewUserApproach(t))
	}
	filter := mock_models.NewUserApproachFilter(t)
	type fields struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserApproachFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.UserApproach
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userApproachRepository.EXPECT().List(ctx, filter).Return(userApproaches, nil)
				userApproachRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    userApproaches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				userApproachRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
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
				userApproachRepository.EXPECT().List(ctx, filter).Return(userApproaches, nil)
				userApproachRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
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
			u := &UserApproachUseCase{
				userApproachRepository: tt.fields.userApproachRepository,
				logger:                 tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserApproachUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserApproachUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	create := mock_models.NewUserApproachCreate(t)
	type fields struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.UserApproachCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userApproachRepository.EXPECT().
					Create(ctx, &models.UserApproach{}).
					Return(nil)
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    &models.UserApproach{},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				userApproachRepository.EXPECT().
					Create(ctx, &models.UserApproach{}).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
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
		//		userApproachRepository: userApproachRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.UserApproachCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserApproachUseCase{
				userApproachRepository: tt.fields.userApproachRepository,
				logger:                 tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	update := mock_models.NewUserApproachUpdate(t)
	type fields struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.UserApproachUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userApproachRepository.EXPECT().
					Get(ctx, update.ID).Return(userApproach, nil)
				userApproachRepository.EXPECT().
					Update(ctx, userApproach).Return(nil)
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				userApproachRepository.EXPECT().
					Get(ctx, update.ID).
					Return(userApproach, nil)
				userApproachRepository.EXPECT().
					Update(ctx, userApproach).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "UserApproach not found",
			setup: func() {
				userApproachRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
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
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx: ctx,
				update: &models.UserApproachUpdate{
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
			u := &UserApproachUseCase{
				userApproachRepository: tt.fields.userApproachRepository,
				logger:                 tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApproachRepository := mock_repositories.NewMockUserApproachRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	type fields struct {
		userApproachRepository repositories.UserApproachRepository
		logger                 log.Logger
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
				userApproachRepository.EXPECT().
					Delete(ctx, userApproach.ID).
					Return(nil)
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx: ctx,
				id:  userApproach.ID,
			},
			wantErr: nil,
		},
		{
			name: "UserApproach not found",
			setup: func() {
				userApproachRepository.EXPECT().
					Delete(ctx, userApproach.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				userApproachRepository: userApproachRepository,
				logger:                 logger,
			},
			args: args{
				ctx: ctx,
				id:  userApproach.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserApproachUseCase{
				userApproachRepository: tt.fields.userApproachRepository,
				logger:                 tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
