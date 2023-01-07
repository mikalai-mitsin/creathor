package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	mock_usecases "github.com/018bf/example/internal/domain/usecases/mock"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

func TestNewUserSessionInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase        usecases.AuthUseCase
		userSessionUseCase usecases.UserSessionUseCase
		logger             log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.UserSessionInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
			},
			want: &UserSessionInterceptor{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserSessionInterceptor(tt.args.userSessionUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserSessionInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	type fields struct {
		authUseCase        usecases.AuthUseCase
		userSessionUseCase usecases.UserSessionUseCase
		logger             log.Logger
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
		want    *models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDetail).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionDetail, userSession).
					Return(nil)
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDetail).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionDetail, userSession).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "UserSession not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDetail).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserSessionInterceptor{
				userSessionUseCase: tt.fields.userSessionUseCase,
				authUseCase:        tt.fields.authUseCase,
				logger:             tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	create := mock_models.NewUserSessionCreate(t)
	type fields struct {
		userSessionUseCase usecases.UserSessionUseCase
		authUseCase        usecases.AuthUseCase
		logger             log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.UserSessionCreate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionCreate, create).
					Return(nil)
				userSessionUseCase.EXPECT().Create(ctx, create).Return(userSession, nil)
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionCreate, create).
					Return(nil)
				userSessionUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
			i := &UserSessionInterceptor{
				userSessionUseCase: tt.fields.userSessionUseCase,
				authUseCase:        tt.fields.authUseCase,
				logger:             tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	update := mock_models.NewUserSessionUpdate(t)
	type fields struct {
		userSessionUseCase usecases.UserSessionUseCase
		authUseCase        usecases.AuthUseCase
		logger             log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.UserSessionUpdate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate, userSession).
					Return(nil)
				userSessionUseCase.EXPECT().Update(ctx, update).Return(userSession, nil)
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate, userSession).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate, userSession).
					Return(nil)
				userSessionUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
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
			i := &UserSessionInterceptor{
				userSessionUseCase: tt.fields.userSessionUseCase,
				authUseCase:        tt.fields.authUseCase,
				logger:             tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	userSession := mock_models.NewUserSession(t)
	type fields struct {
		userSessionUseCase usecases.UserSessionUseCase
		authUseCase        usecases.AuthUseCase
		logger             log.Logger
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDelete).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionDelete, userSession).
					Return(nil)
				userSessionUseCase.EXPECT().
					Delete(ctx, userSession.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "UserSession not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDelete).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDelete).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionDelete, userSession).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDelete).
					Return(nil)
				userSessionUseCase.EXPECT().
					Get(ctx, userSession.ID).
					Return(userSession, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionDelete, userSession).
					Return(nil)
				userSessionUseCase.EXPECT().
					Delete(ctx, userSession.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:        authUseCase,
				userSessionUseCase: userSessionUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				id:          userSession.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserSessionInterceptor{
				userSessionUseCase: tt.fields.userSessionUseCase,
				authUseCase:        tt.fields.authUseCase,
				logger:             tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserSessionInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	userSessionUseCase := mock_usecases.NewMockUserSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewUserSessionFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	userSessions := make([]*models.UserSession, 0, count)
	for i := uint64(0); i < count; i++ {
		userSessions = append(userSessions, mock_models.NewUserSession(t))
	}
	type fields struct {
		userSessionUseCase usecases.UserSessionUseCase
		authUseCase        usecases.AuthUseCase
		logger             log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.UserSessionFilter
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionList, filter).
					Return(nil)
				userSessionUseCase.EXPECT().
					List(ctx, filter).
					Return(userSessions, count, nil)
			},
			fields: fields{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    userSessions,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDUserSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDUserSessionList, filter).
					Return(nil)
				userSessionUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				userSessionUseCase: userSessionUseCase,
				authUseCase:        authUseCase,
				logger:             logger,
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
			i := &UserSessionInterceptor{
				userSessionUseCase: tt.fields.userSessionUseCase,
				authUseCase:        tt.fields.authUseCase,
				logger:             tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserSessionInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
