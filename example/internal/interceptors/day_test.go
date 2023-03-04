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

func TestNewDayInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase usecases.AuthUseCase
		dayUseCase  usecases.DayUseCase
		logger      log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.DayInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			want: &DayInterceptor{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewDayInterceptor(tt.args.dayUseCase, tt.args.logger, tt.args.authUseCase); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		authUseCase usecases.AuthUseCase
		dayUseCase  usecases.DayUseCase
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
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayDetail, day).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          models.UUID(day.ID),
				requestUser: requestUser,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayDetail, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Day not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	create := mock_models.NewDayCreate(t)
	type fields struct {
		dayUseCase  usecases.DayUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.DayCreate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayCreate, create).
					Return(nil)
				dayUseCase.EXPECT().Create(ctx, create).Return(day, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayCreate, create).
					Return(nil)
				dayUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	update := mock_models.NewDayUpdate(t)
	type fields struct {
		dayUseCase  usecases.DayUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.DayUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayUpdate, day).
					Return(nil)
				dayUseCase.EXPECT().Update(ctx, update).Return(day, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayUpdate, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayUpdate, day).
					Return(nil)
				dayUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
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
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		dayUseCase  usecases.DayUseCase
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
					HasPermission(ctx, requestUser, models.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayDelete, day).
					Return(nil)
				dayUseCase.EXPECT().
					Delete(ctx, day.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Day not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayDelete, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayDelete, day).
					Return(nil)
				dayUseCase.EXPECT().
					Delete(ctx, day.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          day.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("DayInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	dayUseCase := mock_usecases.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewDayFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listDays := make([]*models.Day, 0, count)
	for i := uint64(0); i < count; i++ {
		listDays = append(listDays, mock_models.NewDay(t))
	}
	type fields struct {
		dayUseCase  usecases.DayUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.DayFilter
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Day
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayList, filter).
					Return(nil)
				dayUseCase.EXPECT().
					List(ctx, filter).
					Return(listDays, count, nil)
			},
			fields: fields{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    listDays,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				dayUseCase:  dayUseCase,
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
					HasPermission(ctx, requestUser, models.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDDayList, filter).
					Return(nil)
				dayUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				dayUseCase:  dayUseCase,
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
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DayInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
