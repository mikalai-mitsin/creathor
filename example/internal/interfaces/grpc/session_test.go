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

func TestNewSessionServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionInterceptor interceptors.SessionInterceptor
		logger             log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.SessionServiceServer
	}{
		{
			name: "ok",
			args: args{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			want: &SessionServiceServer{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionServiceServer(tt.args.sessionInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSessionServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	// create := mock_models.NewSessionCreate(t)
	session := mock_models.NewSession(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                interceptors.SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(session, nil)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.SessionCreate{},
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.
					EXPECT().
					Create(ctx, gomock.Any(), user).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.SessionCreate{},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	id := uuid.NewString()
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                interceptors.SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionDelete
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
				sessionInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionDelete{
					Id: id,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Delete(ctx, models.UUID(id), user).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionDelete{
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
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	session := mock_models.NewSession(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                interceptors.SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().Get(ctx, session.ID, user).Return(session, nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionGet{
					Id: string(session.ID),
				},
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Get(ctx, session.ID, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionGet{
					Id: string(session.ID),
				},
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	filter := mock_models.NewSessionFilter(t)
	var ids []models.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListSession{
		Items: make([]*examplepb.Session, 0, int(count)),
		Count: count,
	}
	listSessions := make([]*models.Session, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewSession(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listSessions = append(listSessions, a)
		response.Items = append(response.Items, decodeSession(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                interceptors.SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().
					List(ctx, filter, user).
					Return(listSessions, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionFilter{
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
				sessionInterceptor.
					EXPECT().
					List(ctx, filter, user).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionFilter{
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
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_interceptors.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	ctx = context.WithValue(ctx, UserKey, user)
	session := mock_models.NewSession(t)
	update := mock_models.NewSessionUpdate(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                interceptors.SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().Update(ctx, update, user).Return(session, nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeSessionUpdate(update),
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Update(ctx, update, user).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeSessionUpdate(update),
			},
			want:    nil,
			wantErr: decodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func Test_decodeSession(t *testing.T) {
	session := mock_models.NewSession(t)
	result := &examplepb.Session{
		Id:          string(session.ID),
		UpdatedAt:   timestamppb.New(session.UpdatedAt),
		CreatedAt:   timestamppb.New(session.CreatedAt),
		Title:       string(session.Title),
		Description: string(session.Description),
	}
	type args struct {
		session *models.Session
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Session
	}{
		{
			name: "ok",
			args: args{
				session: session,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeSession(tt.args.session); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeSessionFilter(t *testing.T) {
	id := models.UUID(uuid.NewString())
	type args struct {
		input *examplepb.SessionFilter
	}
	tests := []struct {
		name string
		args args
		want *models.SessionFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.SessionFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.SessionFilter{
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
			if got := encodeSessionFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
