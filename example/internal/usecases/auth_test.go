package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	mock_repositories "github.com/018bf/example/internal/domain/repositories/mock"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"

	"github.com/golang/mock/gomock"
)

func TestAuthUseCase_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		authRepository repositories.AuthRepository
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		access models.Token
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
				authRepository.EXPECT().GetSubject(ctx, models.Token("mytoken")).Return(user.ID, nil).Times(1)
				userRepository.EXPECT().Get(ctx, user.ID).Return(user, nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "bad user",
			setup: func() {
				authRepository.EXPECT().
					GetSubject(ctx, models.Token("mytoken")).
					Return("", errs.NewBadToken()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    nil,
			wantErr: errs.NewBadToken(),
		},
		{
			name: "user not found",
			setup: func() {
				authRepository.EXPECT().GetSubject(ctx, models.Token("mytoken")).Return(user.ID, nil).Times(1)
				userRepository.EXPECT().Get(ctx, user.ID).Return(nil, errs.NewEntityNotFound()).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Auth(tt.args.ctx, tt.args.access)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	login := mock_models.NewLogin(t)
	user.Email = login.Email
	pair := mock_models.NewTokenPair(t)
	user.SetPassword(login.Password)
	type fields struct {
		authRepository repositories.AuthRepository
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx   context.Context
		login *models.Login
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userRepository.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil).Times(1)
				authRepository.EXPECT().Create(ctx, user).Return(pair, nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(nil, errs.NewEntityNotFound()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "bad password",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(user, nil).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				login: &models.Login{
					Email:    login.Email,
					Password: "mojParol'",
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidParameter("email or password"),
		},
		{
			name: "bad password",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(user, nil).
					Times(1)
				authRepository.EXPECT().Create(ctx, user).
					Return(nil, errs.NewUnexpectedBehaviorError("system errpr")).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("system errpr"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.CreateToken(tt.args.ctx, tt.args.login)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	pair := mock_models.NewTokenPair(t)
	type fields struct {
		authRepository repositories.AuthRepository
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx     context.Context
		refresh models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().RefreshToken(ctx, models.Token("my_r_token")).Return(pair, nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:     ctx,
				refresh: "my_r_token",
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authRepository.EXPECT().
					RefreshToken(ctx, models.Token("my_r_token")).
					Return(nil, errs.NewBadToken()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:     ctx,
				refresh: "my_r_token",
			},
			want:    nil,
			wantErr: errs.NewBadToken(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.RefreshToken(tt.args.ctx, tt.args.refresh)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	type fields struct {
		authRepository repositories.AuthRepository
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		access models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func()
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().Validate(ctx, models.Token("my_token")).Return(nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				access: "my_token",
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authRepository.EXPECT().
					Validate(ctx, models.Token("my_token")).
					Return(errs.NewUnexpectedBehaviorError("error 345")).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				access: "my_token",
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  map[string]string{"details": "error 345"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			if err := u.ValidateToken(tt.args.ctx, tt.args.access); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authRepository       repositories.AuthRepository
		userRepository       repositories.UserRepository
		permissionRepository repositories.PermissionRepository
		logger               log.Logger
	}
	tests := []struct {
		name string
		args args
		want usecases.AuthUseCase
	}{
		{
			name: "ok",
			args: args{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			want: &AuthUseCase{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthUseCase(tt.args.authRepository, tt.args.userRepository, tt.args.permissionRepository, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_HasPermission(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	permissionRepository := mock_repositories.NewMockPermissionRepository(ctrl)
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	user := mock_models.NewUser(t)
	type fields struct {
		authRepository       repositories.AuthRepository
		userRepository       repositories.UserRepository
		logger               log.Logger
		permissionRepository repositories.PermissionRepository
	}
	type args struct {
		in0        context.Context
		in1        *models.User
		permission models.PermissionID
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
				permissionRepository.EXPECT().
					HasPermission(ctx, models.PermissionIDUserList, user).
					Return(nil)
			},
			fields: fields{
				authRepository:       authRepository,
				permissionRepository: permissionRepository,
				userRepository:       userRepository,
				logger:               nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: models.PermissionIDUserList,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				permissionRepository.EXPECT().
					HasPermission(ctx, models.PermissionIDUserList, user).
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository:       authRepository,
				permissionRepository: permissionRepository,
				userRepository:       userRepository,
				logger:               nil,
			},
			args: args{
				in0:        ctx,
				in1:        user,
				permission: models.PermissionIDUserList,
			},
			wantErr: errs.NewPermissionDenied(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthUseCase{
				authRepository:       tt.fields.authRepository,
				userRepository:       tt.fields.userRepository,
				permissionRepository: tt.fields.permissionRepository,
				logger:               tt.fields.logger,
			}
			tt.setup()
			if err := u.HasPermission(tt.args.in0, tt.args.in1, tt.args.permission); !errors.Is(err, tt.wantErr) {
				t.Errorf("HasPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthUseCase_HasObjectPermission(t *testing.T) {
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	permissionRepository := mock_repositories.NewMockPermissionRepository(ctrl)
	type fields struct {
		authRepository       repositories.AuthRepository
		userRepository       repositories.UserRepository
		permissionRepository repositories.PermissionRepository
		logger               log.Logger
	}
	type args struct {
		in0        context.Context
		user       *models.User
		permission models.PermissionID
		object     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func()
	}{
		{
			name: "error",
			setup: func() {
				permissionRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDUserDetail, user, "user").
					Return(errs.NewPermissionDenied())
			},
			fields: fields{
				authRepository:       authRepository,
				permissionRepository: permissionRepository,
				userRepository:       userRepository,
				logger:               nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDUserDetail,
				object:     "user",
			},
			wantErr: errs.NewPermissionDenied(),
		},
		{
			name: "ok",
			setup: func() {
				permissionRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDUserDetail, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository:       authRepository,
				permissionRepository: permissionRepository,
				userRepository:       userRepository,
				logger:               nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDUserDetail,
				object:     user,
			},
			wantErr: nil,
		},
		{
			name: "ok with user",
			setup: func() {
				permissionRepository.EXPECT().
					HasObjectPermission(ctx, models.PermissionIDUserDelete, user, user).
					Return(nil)
			},
			fields: fields{
				authRepository:       authRepository,
				permissionRepository: permissionRepository,
				userRepository:       userRepository,
				logger:               nil,
			},
			args: args{
				in0:        ctx,
				user:       user,
				permission: models.PermissionIDUserDelete,
				object:     user,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthUseCase{
				authRepository:       tt.fields.authRepository,
				userRepository:       tt.fields.userRepository,
				permissionRepository: tt.fields.permissionRepository,
				logger:               tt.fields.logger,
			}
			tt.setup()
			if err := u.HasObjectPermission(tt.args.in0, tt.args.user, tt.args.permission, tt.args.object); !errors.Is(err, tt.wantErr) {
				t.Errorf("HasObjectPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthUseCase_CreateTokenByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := mock_repositories.NewMockUserRepository(ctrl)
	authRepository := mock_repositories.NewMockAuthRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	tokenPair := mock_models.NewTokenPair(t)
	type fields struct {
		authRepository repositories.AuthRepository
		userRepository repositories.UserRepository
		logger         log.Logger
	}
	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().Create(ctx, user).Return(tokenPair, nil)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:  ctx,
				user: user,
			},
			want:    tokenPair,
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				authRepository.EXPECT().Create(ctx, user).Return(nil, errs.NewUnexpectedBehaviorError("asd"))
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         logger,
			},
			args: args{
				ctx:  ctx,
				user: user,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("asd"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthUseCase{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			tt.setup()
			got, err := u.CreateTokenByUser(tt.args.ctx, tt.args.user)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateTokenByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTokenByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
