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

func TestNewArchServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		archInterceptor interceptors.ArchInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.ArchServiceServer
	}{
		{
			name: "ok",
			args: args{
				archInterceptor: archInterceptor,
				logger:          logger,
			},
			want: &ArchServiceServer{
				archInterceptor: archInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArchServiceServer(tt.args.archInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewArchServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	// create := mock_models.NewArchCreate(t)
	arch := mock_models.NewArch(t)
	type fields struct {
		UnimplementedArchServiceServer examplepb.UnimplementedArchServiceServer
		archInterceptor                interceptors.ArchInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArchCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(arch, nil)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.ArchCreate{},
			},
			want:    decodeArch(arch),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				archInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.ArchCreate{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArchServiceServer{
				UnimplementedArchServiceServer: tt.fields.UnimplementedArchServiceServer,
				archInterceptor:                tt.fields.archInterceptor,
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

func TestArchServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	id := uuid.NewString()
	type fields struct {
		UnimplementedArchServiceServer examplepb.UnimplementedArchServiceServer
		archInterceptor                interceptors.ArchInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArchDelete
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
				archInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchDelete{
					Id: id,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				archInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchDelete{
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
			s := ArchServiceServer{
				UnimplementedArchServiceServer: tt.fields.UnimplementedArchServiceServer,
				archInterceptor:                tt.fields.archInterceptor,
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

func TestArchServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	arch := mock_models.NewArch(t)
	type fields struct {
		UnimplementedArchServiceServer examplepb.UnimplementedArchServiceServer
		archInterceptor                interceptors.ArchInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArchGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archInterceptor.EXPECT().Get(ctx, arch.ID, user).Return(arch, nil).Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchGet{
					Id: string(arch.ID),
				},
			},
			want:    decodeArch(arch),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				archInterceptor.EXPECT().Get(ctx, arch.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchGet{
					Id: string(arch.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := ArchServiceServer{
				UnimplementedArchServiceServer: tt.fields.UnimplementedArchServiceServer,
				archInterceptor:                tt.fields.archInterceptor,
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

func TestArchServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	filter := mock_models.NewArchFilter(t)
	var ids []models.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListArch{
		Items: make([]*examplepb.Arch, 0, int(count)),
		Count: count,
	}
	listArches := make([]*models.Arch, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewArch(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listArches = append(listArches, a)
		response.Items = append(response.Items, decodeArch(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedArchServiceServer examplepb.UnimplementedArchServiceServer
		archInterceptor                interceptors.ArchInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArchFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListArch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archInterceptor.EXPECT().
					List(ctx, filter, user).
					Return(listArches, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchFilter{
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
				archInterceptor.
					EXPECT().
					List(ctx, filter, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArchFilter{
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
			s := ArchServiceServer{
				UnimplementedArchServiceServer: tt.fields.UnimplementedArchServiceServer,
				archInterceptor:                tt.fields.archInterceptor,
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

func TestArchServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archInterceptor := mock_interceptors.NewMockArchInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	arch := mock_models.NewArch(t)
	update := mock_models.NewArchUpdate(t)
	type fields struct {
		UnimplementedArchServiceServer examplepb.UnimplementedArchServiceServer
		archInterceptor                interceptors.ArchInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArchUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archInterceptor.EXPECT().Update(ctx, update, user).Return(arch, nil).Times(1)
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeArchUpdate(update),
			},
			want:    decodeArch(arch),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				archInterceptor.EXPECT().Update(ctx, update, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedArchServiceServer: examplepb.UnimplementedArchServiceServer{},
				archInterceptor:                archInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeArchUpdate(update),
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArchServiceServer{
				UnimplementedArchServiceServer: tt.fields.UnimplementedArchServiceServer,
				archInterceptor:                tt.fields.archInterceptor,
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

func Test_decodeArch(t *testing.T) {
	arch := mock_models.NewArch(t)
	result := &examplepb.Arch{
		Id:          string(arch.ID),
		UpdatedAt:   timestamppb.New(arch.UpdatedAt),
		CreatedAt:   timestamppb.New(arch.CreatedAt),
		Name:        string(arch.Name),
		Title:       string(arch.Title),
		Description: string(arch.Description),
		Tags:        []string{},
		Versions:    []uint64{},
		Release:     timestamppb.New(arch.Release),
		Tested:      timestamppb.New(arch.Tested),
	}
	for _, param := range arch.Tags {
		result.Tags = append(result.Tags, string(param))
	}
	for _, param := range arch.Versions {
		result.Versions = append(result.Versions, uint64(param))
	}
	type args struct {
		arch *models.Arch
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Arch
	}{
		{
			name: "ok",
			args: args{
				arch: arch,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeArch(tt.args.arch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeArchFilter(t *testing.T) {
	id := models.UUID(uuid.NewString())
	type args struct {
		input *examplepb.ArchFilter
	}
	tests := []struct {
		name string
		args args
		want *models.ArchFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.ArchFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.ArchFilter{
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
			if got := encodeArchFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
