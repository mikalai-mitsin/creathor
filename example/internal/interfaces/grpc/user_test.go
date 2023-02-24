package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/utils"
	"github.com/jaswdr/faker"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	mock_interceptors "github.com/018bf/example/internal/domain/interceptors/mock"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewUserServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userInterceptor interceptors.UserInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.UserServiceServer
	}{
		{
			name: "ok",
			args: args{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			want: &UserServiceServer{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserServiceServer(tt.args.userInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUserServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	userID := uuid.NewString()
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                interceptors.UserInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().Delete(ctx, models.UUID(userID), user).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserDelete{
					Id: userID,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Delete(ctx, models.UUID(userID), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserDelete{
					Id: userID,
				},
			},
			want: nil,
			wantErr: decodeError(&errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "i error",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := u.Delete(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                interceptors.UserInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().Get(ctx, user.ID, user).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserGet{
					Id: string(user.ID),
				},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Get(ctx, user.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserGet{
					Id: string(user.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	userFilter := mock_models.NewUserFilter(t)
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.Users{
		Items: make([]*examplepb.User, 0, int(count)),
		Count: count,
	}
	users := make([]*models.User, 0, int(count))
	for i := 0; i < int(count); i++ {
		u := mock_models.NewUser(t)
		users = append(users, u)
		response.Items = append(response.Items, decodeUser(u))
	}
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                interceptors.UserInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Users
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().
					List(ctx, userFilter, user).
					Return(users, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserFilter{
					PageNumber: wrapperspb.UInt64(*userFilter.PageNumber),
					PageSize:   wrapperspb.UInt64(*userFilter.PageSize),
					Search:     wrapperspb.String(*userFilter.Search),
					OrderBy:    userFilter.OrderBy,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().List(ctx, &models.UserFilter{}, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.UserFilter{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := u.List(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	create := mock_models.NewUserCreate(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                interceptors.UserInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.Signup
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().Create(ctx, create, user).Return(user, nil).Times(1)

			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.Signup{
					Email:    create.Email,
					Password: create.Password,
				},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Create(ctx, create, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.Signup{
					Email:    create.Email,
					Password: create.Password,
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := u.Signup(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Signup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_interceptors.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	userUpdate := mock_models.NewUserUpdate(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                interceptors.UserInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userInterceptor.EXPECT().Update(ctx, &models.UserUpdate{
					ID:        userUpdate.ID,
					FirstName: userUpdate.FirstName,
					LastName:  userUpdate.LastName,
					Password:  userUpdate.Password,
					Email:     userUpdate.Email,
				}, user).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserUpdate{
					Id:        string(userUpdate.ID),
					FirstName: wrapperspb.String(*userUpdate.FirstName),
					LastName:  wrapperspb.String(*userUpdate.LastName),
					Email:     wrapperspb.String(*userUpdate.Email),
					Password:  wrapperspb.String(*userUpdate.Password),
				},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Update(ctx, &models.UserUpdate{
					ID:        userUpdate.ID,
					FirstName: userUpdate.FirstName,
					LastName:  userUpdate.LastName,
					Password:  userUpdate.Password,
					Email:     userUpdate.Email,
				}, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserUpdate{
					Id:        string(userUpdate.ID),
					FirstName: wrapperspb.String(*userUpdate.FirstName),
					LastName:  wrapperspb.String(*userUpdate.LastName),
					Email:     wrapperspb.String(*userUpdate.Email),
					Password:  wrapperspb.String(*userUpdate.Password),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeUser(t *testing.T) {
	user := mock_models.NewUser(t)
	type args struct {
		user *models.User
	}
	tests := []struct {
		name string
		args args
		want *examplepb.User
	}{
		{
			name: "ok",
			args: args{
				user: user,
			},
			want: &examplepb.User{
				Id:        string(user.ID),
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				CreatedAt: timestamppb.New(user.CreatedAt),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeUserFilter(t *testing.T) {
	type args struct {
		input *examplepb.UserFilter
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *models.UserFilter
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				input: &examplepb.UserFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
				},
			},
			want: &models.UserFilter{
				PageSize:   utils.Pointer(uint64(5)),
				PageNumber: utils.Pointer(uint64(2)),
				Search:     utils.Pointer("my name is"),
				OrderBy:    []string{"created_at", "id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := encodeUserFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
