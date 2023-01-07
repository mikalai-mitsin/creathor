package interceptors

import (
	"context"
	"errors"
	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/usecases"
	mock_usecases "github.com/018bf/example/internal/domain/usecases/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"syreclabs.com/go/faker"
	"testing"
)

func TestNewUserInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.UserInterceptor
	}{
		{
			name: "ok",
			setup: func() {

			},
			args: args{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			want: &UserInterceptor{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserInterceptor(tt.args.userUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
		in2 *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDetail).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserDetail, user).Return(nil)

			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDetail).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "user not found",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDetail).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(nil, errs.NewEntityNotFound())

			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDetail).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserDetail, user).Return(errs.NewPermissionDeniedError())

			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.in2)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	create := mock_models.NewUserCreate(t)
	type fields struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.UserCreate
		in2    *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserCreate).Return(nil)
				userUseCase.EXPECT().Create(ctx, create).Return(user, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				in2:    user,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserCreate).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				in2:    user,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserCreate).Return(nil)
				userUseCase.EXPECT().Create(ctx, create).Return(user, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
				in2:    user,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.in2)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	update := mock_models.NewUserUpdate(t)
	type fields struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.UserUpdate
		in2    *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserUpdate).Return(nil)
				userUseCase.EXPECT().Get(ctx, update.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserUpdate, user).Return(nil)
				userUseCase.EXPECT().Update(ctx, update).Return(user, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				in2:    user,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserUpdate).Return(nil)
				userUseCase.EXPECT().Get(ctx, update.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserUpdate, user).Return(nil)
				userUseCase.EXPECT().Update(ctx, update).Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				in2:    user,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserUpdate).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				in2:    user,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "user not found",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserUpdate).Return(nil)
				userUseCase.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())

			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				in2:    user,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserUpdate).Return(nil)
				userUseCase.EXPECT().Get(ctx, update.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserUpdate, user).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
				in2:    user,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.in2)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
		in2 *models.User
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
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDelete).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserDelete, user).Return(nil)
				userUseCase.EXPECT().Delete(ctx, user.ID).Return(nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDelete).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserDelete, user).Return(nil)
				userUseCase.EXPECT().Delete(ctx, user.ID).Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDelete).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "user not found",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDelete).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserDelete).Return(nil)
				userUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil)
				authUseCase.EXPECT().HasObjectPermission(ctx, user, models.PermissionIDUserDelete, user).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
				in2: user,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.in2); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userUseCase := mock_usecases.NewMockUserUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	filter := mock_models.NewUserFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	users := make([]*models.User, 0, count)
	for i := uint64(0); i < count; i++ {
		users = append(users, mock_models.NewUser(t))
	}
	type fields struct {
		userUseCase usecases.UserUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserFilter
		in2    *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.User
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserList).Return(nil)
				userUseCase.EXPECT().List(ctx, filter).Return(users, count, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				in2:    user,
			},
			want:    users,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserList).Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				in2:    user,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().HasPermission(ctx, user, models.PermissionIDUserList).Return(nil)
				userUseCase.EXPECT().List(ctx, filter).Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				userUseCase: userUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
				in2:    user,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.in2)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
