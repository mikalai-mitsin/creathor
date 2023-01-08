package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/interfaces/postgres"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"reflect"
	"syreclabs.com/go/faker"
	"testing"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/log"
	"github.com/jmoiron/sqlx"
)

func TestNewUserSessionRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		database *sqlx.DB
		logger   log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  repositories.UserSessionRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &UserSessionRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserSessionRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserSessionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.user_sessions"
	userSession := mock_models.NewUserSession(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.UserSession
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
				// TODO: add args
				mock.ExpectQuery(query).
					WithArgs().
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(userSession.ID, userSession.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				// TODO: add args
				mock.ExpectQuery(query).
					WithArgs().
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserSessionRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT user_sessions.id, user_sessions.updated_at, user_sessions.created_at FROM public.user_sessions WHERE id = \\$1 LIMIT 1"
	userSession := mock_models.NewUserSession(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
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
		want    *models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newUserSessionRows(t, []*models.UserSession{userSession})
				mock.ExpectQuery(query).WithArgs(userSession.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  userSession.ID,
			},
			want:    userSession,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(userSession.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_session_id", userSession.ID),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(userSession.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound().WithParam("user_session_id", userSession.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionRepository_List(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var userSessions []*models.UserSession
	for i := 0; i < faker.RandomInt(1, 20); i++ {
		userSessions = append(userSessions, mock_models.NewUserSession(t))
	}
	filter := mock_models.NewUserSessionFilter(t)
	query := "SELECT user_sessions.id, user_sessions.updated_at, user_sessions.created_at FROM public.user_sessions"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserSessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.UserSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newUserSessionRows(t, userSessions))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    userSessions,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error": "test error",
				},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSessionRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSessionRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	userSession := mock_models.NewUserSession(t)
	query := `UPDATE public.user_sessions`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.UserSession
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
				// TODO: set args
				mock.ExpectExec(query).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				// TODO: set args
				mock.ExpectExec(query).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: errs.NewEntityNotFound().WithParam("user_session_id", userSession.ID),
		},
		{
			name: "database error",
			setup: func() {
				// TODO: set args
				mock.ExpectExec(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_session_id", userSession.ID),
		},
		{
			name: "unexpected error",
			setup: func() {
				// TODO: set args
				mock.ExpectExec(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_session_id", userSession.ID),
		},
		{
			name: "result error",
			setup: func() {
				// TODO: set args
				mock.ExpectExec(query).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userSession,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_session_id", userSession.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserSessionRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	userSession := mock_models.NewUserSession(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
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
			fields: fields{
				database: db,
				logger:   logger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_sessions WHERE id = \\$1").
					WithArgs(userSession.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_sessions WHERE id = \\$1").
					WithArgs(userSession.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			wantErr: errs.NewEntityNotFound().WithParam("user_session_id", userSession.ID),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_sessions WHERE id = \\$1").
					WithArgs(userSession.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_session_id", userSession.ID),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_sessions WHERE id = \\$1").
					WithArgs(userSession.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userSession.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_session_id", userSession.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserSessionRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserSessionRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := `SELECT count\(id\) FROM public.user_sessions`
	ctx := context.Background()
	filter := mock_models.NewUserSessionFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserSessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(1))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "bad return type",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow("one"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: 0,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error": "sql: Scan error on column index 0, name \"count\": converting driver.Value type string (\"one\") to a uint64: invalid syntax",
				},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    0,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserSessionRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Count() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newUserSessionRows(t *testing.T, userSessions []*models.UserSession) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		// TODO: add columns
		"updated_at",
		"created_at",
	})
	for _, userSession := range userSessions {
		rows.AddRow(
			userSession.ID,
			// TODO: add values
			userSession.UpdatedAt,
			userSession.CreatedAt,
		)
	}
	return rows
}
