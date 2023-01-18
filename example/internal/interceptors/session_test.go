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

func TestNewSessionInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionUseCase usecases.SessionUseCase
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.SessionInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			want: &SessionInterceptor{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewSessionInterceptor(tt.args.sessionUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		sessionUseCase usecases.SessionUseCase
		logger         log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "Session not found",
			setup: func() {
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	create := mock_models.NewSessionCreate(t)
	type fields struct {
		sessionUseCase usecases.SessionUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.SessionCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionUseCase.EXPECT().Create(ctx, create).Return(session, nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				sessionUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	update := mock_models.NewSessionUpdate(t)
	type fields struct {
		sessionUseCase usecases.SessionUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.SessionUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionUseCase.EXPECT().Update(ctx, update).Return(session, nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				sessionUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		sessionUseCase usecases.SessionUseCase
		logger         log.Logger
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
				sessionUseCase.EXPECT().
					Delete(ctx, session.ID).
					Return(nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				sessionUseCase.EXPECT().
					Delete(ctx, session.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				logger:         tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSessionInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCase := mock_usecases.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewSessionFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	sessions := make([]*models.Session, 0, count)
	for i := uint64(0); i < count; i++ {
		sessions = append(sessions, mock_models.NewSession(t))
	}
	type fields struct {
		sessionUseCase usecases.SessionUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.SessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Session
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionUseCase.EXPECT().
					List(ctx, filter).
					Return(sessions, count, nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    sessions,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				sessionUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				logger:         tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SessionInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
