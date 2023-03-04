package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	mock_repositories "github.com/018bf/example/internal/domain/repositories/mock"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	mock_clock "github.com/018bf/example/pkg/clock/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
)

func TestNewSessionUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionRepository repositories.SessionRepository
		clock             clock.Clock
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
				clock:             clockMock,
				logger:            logger,
			},
			want: &SessionUseCase{
				sessionRepository: sessionRepository,
				clock:             clockMock,
				logger:            logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewSessionUseCase(tt.args.sessionRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
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
		id  models.UUID
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
				sessionRepository.EXPECT().
					Get(ctx, session.ID).
					Return(nil, errs.NewEntityNotFound())
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
	var listSessions []*models.Session
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listSessions = append(listSessions, mock_models.NewSession(t))
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
				sessionRepository.EXPECT().List(ctx, filter).Return(listSessions, nil)
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
			want:    listSessions,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				sessionRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
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
				sessionRepository.EXPECT().List(ctx, filter).Return(listSessions, nil)
				sessionRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
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
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewSessionCreate(t)
	now := time.Now().UTC()
	type fields struct {
		sessionRepository repositories.SessionRepository
		clock             clock.Clock
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
				clockMock.EXPECT().Now().Return(now)
				sessionRepository.EXPECT().
					Create(
						ctx,
						&models.Session{
							Title:       create.Title,
							Description: create.Description,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Session{
				ID:          "",
				Title:       create.Title,
				Description: create.Description,
				UpdatedAt:   now,
				CreatedAt:   now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				sessionRepository.EXPECT().
					Create(
						ctx,
						&models.Session{
							ID:          "",
							Title:       create.Title,
							Description: create.Description,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
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
				ctx:    ctx,
				create: &models.SessionCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(map[string]string{
				"title":       "cannot be blank",
				"description": "cannot be blank",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &SessionUseCase{
				sessionRepository: tt.fields.sessionRepository,
				clock:             tt.fields.clock,
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
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewSessionUpdate(t)
	now := session.UpdatedAt
	type fields struct {
		sessionRepository repositories.SessionRepository
		clock             clock.Clock
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
				clockMock.EXPECT().Now().Return(now)
				sessionRepository.EXPECT().
					Get(ctx, update.ID).Return(session, nil)
				sessionRepository.EXPECT().
					Update(ctx, session).Return(nil)
			},
			fields: fields{
				sessionRepository: sessionRepository,
				clock:             clockMock,
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
				clockMock.EXPECT().Now().Return(now)
				sessionRepository.EXPECT().
					Get(ctx, update.ID).
					Return(session, nil)
				sessionRepository.EXPECT().
					Update(ctx, session).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				sessionRepository: sessionRepository,
				clock:             clockMock,
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
				clock:             clockMock,
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
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				update: &models.SessionUpdate{
					ID: models.UUID("baduuid"),
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
				clock:             tt.fields.clock,
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
		id  models.UUID
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
