package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	mock_interceptors "github.com/018bf/example/internal/domain/interceptors/mock"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/018bf/example/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewPlanServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		planInterceptor interceptors.PlanInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.PlanServiceServer
	}{
		{
			name: "ok",
			args: args{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			want: &PlanServiceServer{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlanServiceServer(tt.args.planInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	// create := mock_models.NewPlanCreate(t)
	plan := mock_models.NewPlan(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                interceptors.PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(plan, nil)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PlanCreate{},
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PlanCreate{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	id := uuid.NewString()
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                interceptors.PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanDelete
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
				planInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanDelete{
					Id: id,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanDelete{
					Id: id,
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
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
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

func TestPlanServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	plan := mock_models.NewPlan(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                interceptors.PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().Get(ctx, plan.ID, user).Return(plan, nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanGet{
					Id: string(plan.ID),
				},
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Get(ctx, plan.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanGet{
					Id: string(plan.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
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

func TestPlanServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	filter := mock_models.NewPlanFilter(t)
	var ids []models.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListPlan{
		Items: make([]*examplepb.Plan, 0, int(count)),
		Count: count,
	}
	listPlans := make([]*models.Plan, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewPlan(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listPlans = append(listPlans, a)
		response.Items = append(response.Items, decodePlan(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                interceptors.PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListPlan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().
					List(ctx, filter, user).
					Return(listPlans, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.
					EXPECT().
					List(ctx, filter, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
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

func TestPlanServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_interceptors.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	plan := mock_models.NewPlan(t)
	update := mock_models.NewPlanUpdate(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                interceptors.PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().Update(ctx, gomock.Any(), user).Return(plan, nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodePlanUpdate(update),
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Update(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodePlanUpdate(update),
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
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

func Test_decodePlan(t *testing.T) {
	plan := mock_models.NewPlan(t)
	result := &examplepb.Plan{
		Id:          string(plan.ID),
		UpdatedAt:   timestamppb.New(plan.UpdatedAt),
		CreatedAt:   timestamppb.New(plan.CreatedAt),
		Name:        string(plan.Name),
		Repeat:      uint64(plan.Repeat),
		EquipmentId: string(plan.EquipmentID),
	}
	type args struct {
		plan *models.Plan
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Plan
	}{
		{
			name: "ok",
			args: args{
				plan: plan,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodePlan(tt.args.plan); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodePlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodePlanFilter(t *testing.T) {
	id := models.UUID(uuid.NewString())
	type args struct {
		input *examplepb.PlanFilter
	}
	tests := []struct {
		name string
		args args
		want *models.PlanFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.PlanFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.PlanFilter{
				PageSize:   utils.Pointer(uint64(5)),
				PageNumber: utils.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				Search:     utils.Pointer("my name is"),
				IDs:        []models.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodePlanFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
