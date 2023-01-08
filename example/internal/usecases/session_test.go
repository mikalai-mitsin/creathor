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
	"syreclabs.com/go/faker"
)

func TestNewSessionUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.SessionUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			want: &SessionUseCase{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewSessionUseCase(tt.args.sessionRepository, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
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
				sessionRepository.EXPECT().Get(ctx, session.ID).Return(session, nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
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
				sessionRepository.EXPECT().Get(ctx, session.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
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
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var sessions []*models.Session
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		sessions = append(sessions, mock_models.NewSession(t))
	}
	filter := mock_models.NewSessionFilter(t)
	type fields struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
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
				sessionRepository.EXPECT().List(ctx, filter).Return(sessions, nil)
				sessionRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
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
				sessionRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "count error",
			setup: func() {
				sessionRepository.EXPECT().List(ctx, filter).Return(sessions, nil)
				sessionRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				logger:            tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SessionUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSessionUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	create := mock_models.NewSessionCreate(t)
	type fields struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
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
				sessionRepository.EXPECT().
					Create(ctx, &models.Session{}).
					Return(nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    &models.Session{},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				sessionRepository.EXPECT().
					Create(ctx, &models.Session{}).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		// TODO: Add validation rules or delete this case
		//{
		//	name: "invalid",
		//	setup: func() {
		//	},
		//	fields: fields{
		//		sessionRepository: sessionRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.SessionCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	update := mock_models.NewSessionUpdate(t)
	type fields struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.SessionUpdate
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
				sessionRepository.EXPECT().
					Get(ctx, update.ID).Return(session, nil)
				sessionRepository.EXPECT().
					Update(ctx, session).Return(nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
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
				sessionRepository.EXPECT().
					Get(ctx, update.ID).
					Return(session, nil)
				sessionRepository.EXPECT().
					Update(ctx, session).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Session not found",
			setup: func() {
				sessionRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				update: &models.SessionUpdate{
					ID: faker.Number().Number(1),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		sessionRepository repositories.SessionRepository
		logger            log.Logger
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionRepository.EXPECT().
					Delete(ctx, session.ID).
					Return(nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: nil,
		},
		{
			name: "Session not found",
			setup: func() {
				sessionRepository.EXPECT().
					Delete(ctx, session.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				sessionRepository: sessionRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				logger:            tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
