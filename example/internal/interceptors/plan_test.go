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
	"github.com/jaswdr/faker"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

func TestNewPlanInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase usecases.AuthUseCase
		planUseCase usecases.PlanUseCase
		logger      log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.PlanInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			want: &PlanInterceptor{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPlanInterceptor(tt.args.planUseCase, tt.args.logger, tt.args.authUseCase); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	type fields struct {
		authUseCase usecases.AuthUseCase
		planUseCase usecases.PlanUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		id          models.UUID
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDetail).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDetail, plan).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          models.UUID(plan.ID),
				requestUser: requestUser,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDetail).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDetail, plan).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Plan not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDetail).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PlanInterceptor{
				planUseCase: tt.fields.planUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	create := mock_models.NewPlanCreate(t)
	type fields struct {
		planUseCase usecases.PlanUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.PlanCreate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanCreate, create).
					Return(nil)
				planUseCase.EXPECT().Create(ctx, create).Return(plan, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanCreate, create).
					Return(nil)
				planUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
			i := &PlanInterceptor{
				planUseCase: tt.fields.planUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	update := mock_models.NewPlanUpdate(t)
	type fields struct {
		planUseCase usecases.PlanUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.PlanUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanUpdate, plan).
					Return(nil)
				planUseCase.EXPECT().Update(ctx, update).Return(plan, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanUpdate, plan).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanUpdate, plan).
					Return(nil)
				planUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
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
			i := &PlanInterceptor{
				planUseCase: tt.fields.planUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	type fields struct {
		planUseCase usecases.PlanUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		id          models.UUID
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanDelete).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDelete, plan).
					Return(nil)
				planUseCase.EXPECT().
					Delete(ctx, plan.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Plan not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDelete).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDelete).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDelete, plan).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDelete).
					Return(nil)
				planUseCase.EXPECT().
					Get(ctx, plan.ID).
					Return(plan, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDelete, plan).
					Return(nil)
				planUseCase.EXPECT().
					Delete(ctx, plan.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				planUseCase: planUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          plan.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PlanInterceptor{
				planUseCase: tt.fields.planUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("PlanInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	planUseCase := mock_usecases.NewMockPlanUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewPlanFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listPlans := make([]*models.Plan, 0, count)
	for i := uint64(0); i < count; i++ {
		listPlans = append(listPlans, mock_models.NewPlan(t))
	}
	type fields struct {
		planUseCase usecases.PlanUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.PlanFilter
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Plan
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanList, filter).
					Return(nil)
				planUseCase.EXPECT().
					List(ctx, filter).
					Return(listPlans, count, nil)
			},
			fields: fields{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    listPlans,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPlanList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDPlanList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPlanList, filter).
					Return(nil)
				planUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				planUseCase: planUseCase,
				authUseCase: authUseCase,
				logger:      logger,
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
			i := &PlanInterceptor{
				planUseCase: tt.fields.planUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlanInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
