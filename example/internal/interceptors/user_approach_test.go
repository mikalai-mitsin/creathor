package interceptors

import (
	"context"
	"errors"
	"github.com/018bf/creathor/internal/domain/errs"
	mock_models "github.com/018bf/creathor/internal/domain/models/mock"
	mock_usecases "github.com/018bf/creathor/internal/domain/usecases/mock"
	mock_log "github.com/018bf/creathor/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"syreclabs.com/go/faker"
	"testing"

	"github.com/018bf/creathor/internal/domain/interceptors"
	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/internal/domain/usecases"
	"github.com/018bf/creathor/pkg/log"
)

func TestNewUserApproachInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase         usecases.AuthUseCase
		userApproachUseCase usecases.UserApproachUseCase
		logger              log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.UserApproachInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
			want: &UserApproachInterceptor{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserApproachInterceptor(tt.args.userApproachUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserApproachInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	type fields struct {
		authUseCase         usecases.AuthUseCase
		userApproachUseCase usecases.UserApproachUseCase
		logger              log.Logger
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
		want    *models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDetail).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDetail, userApproach).
					Return(nil)
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDetail).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDetail, userApproach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "UserApproach not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDetail).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserApproachInterceptor{
				userApproachUseCase: tt.fields.userApproachUseCase,
				authUseCase:         tt.fields.authUseCase,
				logger:              tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	create := mock_models.NewUserApproachCreate(t)
	type fields struct {
		userApproachUseCase usecases.UserApproachUseCase
		authUseCase         usecases.AuthUseCase
		logger              log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.UserApproachCreate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachCreate, create).
					Return(nil)
				userApproachUseCase.EXPECT().Create(ctx, create).Return(userApproach, nil)
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachCreate, create).
					Return(nil)
				userApproachUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserApproachInterceptor{
				userApproachUseCase: tt.fields.userApproachUseCase,
				authUseCase:         tt.fields.authUseCase,
				logger:              tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	update := mock_models.NewUserApproachUpdate(t)
	type fields struct {
		userApproachUseCase usecases.UserApproachUseCase
		authUseCase         usecases.AuthUseCase
		logger              log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.UserApproachUpdate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate, userApproach).
					Return(nil)
				userApproachUseCase.EXPECT().Update(ctx, update).Return(userApproach, nil)
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate, userApproach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate, userApproach).
					Return(nil)
				userApproachUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserApproachInterceptor{
				userApproachUseCase: tt.fields.userApproachUseCase,
				authUseCase:         tt.fields.authUseCase,
				logger:              tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userApproach := mock_models.NewUserApproach(t)
	type fields struct {
		userApproachUseCase usecases.UserApproachUseCase
		authUseCase         usecases.AuthUseCase
		logger              log.Logger
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDelete, userApproach).
					Return(nil)
				userApproachUseCase.EXPECT().
					Delete(ctx, userApproach.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "UserApproach not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDelete, userApproach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete).
					Return(nil)
				userApproachUseCase.EXPECT().
					Get(ctx, userApproach.ID).
					Return(userApproach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDelete, userApproach).
					Return(nil)
				userApproachUseCase.EXPECT().
					Delete(ctx, userApproach.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:         authUseCase,
				userApproachUseCase: userApproachUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				id:          userApproach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserApproachInterceptor{
				userApproachUseCase: tt.fields.userApproachUseCase,
				authUseCase:         tt.fields.authUseCase,
				logger:              tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserApproachInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userApproachUseCase := mock_usecases.NewMockUserApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewUserApproachFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	userApproaches := make([]*models.UserApproach, 0, count)
	for i := uint64(0); i < count; i++ {
		userApproaches = append(userApproaches, mock_models.NewUserApproach(t))
	}
	type fields struct {
		userApproachUseCase usecases.UserApproachUseCase
		authUseCase         usecases.AuthUseCase
		logger              log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.UserApproachFilter
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachList, filter).
					Return(nil)
				userApproachUseCase.EXPECT().
					List(ctx, filter).
					Return(userApproaches, count, nil)
			},
			fields: fields{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    userApproaches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachList, filter).
					Return(nil)
				userApproachUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				userApproachUseCase: userApproachUseCase,
				authUseCase:         authUseCase,
				logger:              logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserApproachInterceptor{
				userApproachUseCase: tt.fields.userApproachUseCase,
				authUseCase:         tt.fields.authUseCase,
				logger:              tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserApproachInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
