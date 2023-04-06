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

func TestNewEquipmentServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		equipmentInterceptor interceptors.EquipmentInterceptor
		logger               log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.EquipmentServiceServer
	}{
		{
			name: "ok",
			args: args{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
			want: &EquipmentServiceServer{
				equipmentInterceptor: equipmentInterceptor,
				logger:               logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEquipmentServiceServer(tt.args.equipmentInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewEquipmentServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	// create := mock_models.NewEquipmentCreate(t)
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		UnimplementedEquipmentServiceServer examplepb.UnimplementedEquipmentServiceServer
		equipmentInterceptor                interceptors.EquipmentInterceptor
		logger                              log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.EquipmentCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(equipment, nil)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.EquipmentCreate{},
			},
			want:    decodeEquipment(equipment),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				equipmentInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.EquipmentCreate{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := EquipmentServiceServer{
				UnimplementedEquipmentServiceServer: tt.fields.UnimplementedEquipmentServiceServer,
				equipmentInterceptor:                tt.fields.equipmentInterceptor,
				logger:                              tt.fields.logger,
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

func TestEquipmentServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	id := uuid.NewString()
	type fields struct {
		UnimplementedEquipmentServiceServer examplepb.UnimplementedEquipmentServiceServer
		equipmentInterceptor                interceptors.EquipmentInterceptor
		logger                              log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.EquipmentDelete
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
				equipmentInterceptor.EXPECT().
					Delete(ctx, models.UUID(id), user).
					Return(nil).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentDelete{
					Id: id,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				equipmentInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentDelete{
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
			s := EquipmentServiceServer{
				UnimplementedEquipmentServiceServer: tt.fields.UnimplementedEquipmentServiceServer,
				equipmentInterceptor:                tt.fields.equipmentInterceptor,
				logger:                              tt.fields.logger,
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

func TestEquipmentServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		UnimplementedEquipmentServiceServer examplepb.UnimplementedEquipmentServiceServer
		equipmentInterceptor                interceptors.EquipmentInterceptor
		logger                              log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.EquipmentGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentInterceptor.EXPECT().
					Get(ctx, equipment.ID, user).
					Return(equipment, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentGet{
					Id: string(equipment.ID),
				},
			},
			want:    decodeEquipment(equipment),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				equipmentInterceptor.EXPECT().Get(ctx, equipment.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentGet{
					Id: string(equipment.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := EquipmentServiceServer{
				UnimplementedEquipmentServiceServer: tt.fields.UnimplementedEquipmentServiceServer,
				equipmentInterceptor:                tt.fields.equipmentInterceptor,
				logger:                              tt.fields.logger,
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

func TestEquipmentServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	filter := mock_models.NewEquipmentFilter(t)
	var ids []models.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListEquipment{
		Items: make([]*examplepb.Equipment, 0, int(count)),
		Count: count,
	}
	listEquipment := make([]*models.Equipment, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewEquipment(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listEquipment = append(listEquipment, a)
		response.Items = append(response.Items, decodeEquipment(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedEquipmentServiceServer examplepb.UnimplementedEquipmentServiceServer
		equipmentInterceptor                interceptors.EquipmentInterceptor
		logger                              log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.EquipmentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListEquipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentInterceptor.EXPECT().
					List(ctx, filter, user).
					Return(listEquipment, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentFilter{
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
				equipmentInterceptor.
					EXPECT().
					List(ctx, filter, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.EquipmentFilter{
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
			s := EquipmentServiceServer{
				UnimplementedEquipmentServiceServer: tt.fields.UnimplementedEquipmentServiceServer,
				equipmentInterceptor:                tt.fields.equipmentInterceptor,
				logger:                              tt.fields.logger,
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

func TestEquipmentServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	equipmentInterceptor := mock_interceptors.NewMockEquipmentInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	equipment := mock_models.NewEquipment(t)
	update := mock_models.NewEquipmentUpdate(t)
	type fields struct {
		UnimplementedEquipmentServiceServer examplepb.UnimplementedEquipmentServiceServer
		equipmentInterceptor                interceptors.EquipmentInterceptor
		logger                              log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.EquipmentUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				equipmentInterceptor.EXPECT().
					Update(ctx, gomock.Any(), user).
					Return(equipment, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeEquipmentUpdate(update),
			},
			want:    decodeEquipment(equipment),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				equipmentInterceptor.EXPECT().Update(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedEquipmentServiceServer: examplepb.UnimplementedEquipmentServiceServer{},
				equipmentInterceptor:                equipmentInterceptor,
				logger:                              logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeEquipmentUpdate(update),
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := EquipmentServiceServer{
				UnimplementedEquipmentServiceServer: tt.fields.UnimplementedEquipmentServiceServer,
				equipmentInterceptor:                tt.fields.equipmentInterceptor,
				logger:                              tt.fields.logger,
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

func Test_decodeEquipment(t *testing.T) {
	equipment := mock_models.NewEquipment(t)
	result := &examplepb.Equipment{
		Id:        string(equipment.ID),
		UpdatedAt: timestamppb.New(equipment.UpdatedAt),
		CreatedAt: timestamppb.New(equipment.CreatedAt),
		Name:      string(equipment.Name),
		Repeat:    int32(equipment.Repeat),
		Weight:    int32(equipment.Weight),
	}
	type args struct {
		equipment *models.Equipment
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Equipment
	}{
		{
			name: "ok",
			args: args{
				equipment: equipment,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeEquipment(tt.args.equipment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeEquipment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeEquipmentFilter(t *testing.T) {
	id := models.UUID(uuid.NewString())
	type args struct {
		input *examplepb.EquipmentFilter
	}
	tests := []struct {
		name string
		args args
		want *models.EquipmentFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.EquipmentFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.EquipmentFilter{
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
			if got := encodeEquipmentFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
