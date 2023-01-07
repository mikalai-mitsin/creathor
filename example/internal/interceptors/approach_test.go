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

func TestNewApproachInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase     usecases.AuthUseCase
		approachUseCase usecases.ApproachUseCase
		logger          log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.ApproachInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
			},
			want: &ApproachInterceptor{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewApproachInterceptor(tt.args.approachUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApproachInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	type fields struct {
		authUseCase     usecases.AuthUseCase
		approachUseCase usecases.ApproachUseCase
		logger          log.Logger
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
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDetail).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachDetail, approach).
					Return(nil)
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDetail).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachDetail, approach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Approach not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDetail).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ApproachInterceptor{
				approachUseCase: tt.fields.approachUseCase,
				authUseCase:     tt.fields.authUseCase,
				logger:          tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	create := mock_models.NewApproachCreate(t)
	type fields struct {
		approachUseCase usecases.ApproachUseCase
		authUseCase     usecases.AuthUseCase
		logger          log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.ApproachCreate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachCreate, create).
					Return(nil)
				approachUseCase.EXPECT().Create(ctx, create).Return(approach, nil)
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachCreate, create).
					Return(nil)
				approachUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
			i := &ApproachInterceptor{
				approachUseCase: tt.fields.approachUseCase,
				authUseCase:     tt.fields.authUseCase,
				logger:          tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	update := mock_models.NewApproachUpdate(t)
	type fields struct {
		approachUseCase usecases.ApproachUseCase
		authUseCase     usecases.AuthUseCase
		logger          log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.ApproachUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachUpdate).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachUpdate, approach).
					Return(nil)
				approachUseCase.EXPECT().Update(ctx, update).Return(approach, nil)
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachUpdate).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachUpdate, approach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachUpdate).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachUpdate).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachUpdate, approach).
					Return(nil)
				approachUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
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
			i := &ApproachInterceptor{
				approachUseCase: tt.fields.approachUseCase,
				authUseCase:     tt.fields.authUseCase,
				logger:          tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproachInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	approach := mock_models.NewApproach(t)
	type fields struct {
		approachUseCase usecases.ApproachUseCase
		authUseCase     usecases.AuthUseCase
		logger          log.Logger
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachDelete).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachDelete, approach).
					Return(nil)
				approachUseCase.EXPECT().
					Delete(ctx, approach.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Approach not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDelete).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDelete).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachDelete, approach).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDelete).
					Return(nil)
				approachUseCase.EXPECT().
					Get(ctx, approach.ID).
					Return(approach, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachDelete, approach).
					Return(nil)
				approachUseCase.EXPECT().
					Delete(ctx, approach.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:     authUseCase,
				approachUseCase: approachUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				id:          approach.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ApproachInterceptor{
				approachUseCase: tt.fields.approachUseCase,
				authUseCase:     tt.fields.authUseCase,
				logger:          tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApproachInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	approachUseCase := mock_usecases.NewMockApproachUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewApproachFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	approaches := make([]*models.Approach, 0, count)
	for i := uint64(0); i < count; i++ {
		approaches = append(approaches, mock_models.NewApproach(t))
	}
	type fields struct {
		approachUseCase usecases.ApproachUseCase
		authUseCase     usecases.AuthUseCase
		logger          log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.ApproachFilter
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Approach
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachList, filter).
					Return(nil)
				approachUseCase.EXPECT().
					List(ctx, filter).
					Return(approaches, count, nil)
			},
			fields: fields{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    approaches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDApproachList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDApproachList, filter).
					Return(nil)
				approachUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				approachUseCase: approachUseCase,
				authUseCase:     authUseCase,
				logger:          logger,
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
			i := &ApproachInterceptor{
				approachUseCase: tt.fields.approachUseCase,
				authUseCase:     tt.fields.authUseCase,
				logger:          tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApproachInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproachInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ApproachInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
