package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/018bf/creathor/internal/domain/errs"
	mock_models "github.com/018bf/creathor/internal/domain/models/mock"
	"github.com/018bf/creathor/internal/interfaces/postgres"
	mock_log "github.com/018bf/creathor/pkg/log/mock"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"reflect"
	"syreclabs.com/go/faker"
	"testing"

	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/internal/domain/repositories"
	"github.com/018bf/creathor/pkg/log"
	"github.com/jmoiron/sqlx"
)

func TestNewUserApproachRepository(t *testing.T) {
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
		want  repositories.UserApproachRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &UserApproachRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserApproachRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserApproachRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.user_approaches"
	userApproach := mock_models.NewUserApproach(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.UserApproach
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
						AddRow(userApproach.ID, userApproach.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: userApproach,
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
				card: userApproach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserApproachRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT user_approaches.id, user_approaches.updated_at, user_approaches.created_at FROM public.user_approaches WHERE id = \\$1 LIMIT 1"
	userApproach := mock_models.NewUserApproach(t)
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
		want    *models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newUserApproachRows(t, []*models.UserApproach{userApproach})
				mock.ExpectQuery(query).WithArgs(userApproach.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  userApproach.ID,
			},
			want:    userApproach,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(userApproach.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_approach_id", userApproach.ID),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(userApproach.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound().WithParam("user_approach_id", userApproach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachRepository_List(t *testing.T) {
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
	var userApproaches []*models.UserApproach
	for i := 0; i < faker.RandomInt(1, 20); i++ {
		userApproaches = append(userApproaches, mock_models.NewUserApproach(t))
	}
	filter := mock_models.NewUserApproachFilter(t)
	query := "SELECT user_approaches.id, user_approaches.updated_at, user_approaches.created_at FROM public.user_approaches"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserApproachFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.UserApproach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newUserApproachRows(t, userApproaches))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    userApproaches,
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
			r := &UserApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserApproachRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserApproachRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	userApproach := mock_models.NewUserApproach(t)
	query := `UPDATE public.user_approaches`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.UserApproach
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
				card: userApproach,
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
				card: userApproach,
			},
			wantErr: errs.NewEntityNotFound().WithParam("user_approach_id", userApproach.ID),
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
				card: userApproach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_approach_id", userApproach.ID),
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
				card: userApproach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_approach_id", userApproach.ID),
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
				card: userApproach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_approach_id", userApproach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserApproachRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	userApproach := mock_models.NewUserApproach(t)
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
				mock.ExpectExec("DELETE FROM public.user_approaches WHERE id = \\$1").
					WithArgs(userApproach.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_approaches WHERE id = \\$1").
					WithArgs(userApproach.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			wantErr: errs.NewEntityNotFound().WithParam("user_approach_id", userApproach.ID),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_approaches WHERE id = \\$1").
					WithArgs(userApproach.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_approach_id", userApproach.ID),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.user_approaches WHERE id = \\$1").
					WithArgs(userApproach.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  userApproach.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("user_approach_id", userApproach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserApproachRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserApproachRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := `SELECT count\(id\) FROM public.user_approaches`
	ctx := context.Background()
	filter := mock_models.NewUserApproachFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserApproachFilter
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
			r := &UserApproachRepository{
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

func newUserApproachRows(t *testing.T, userApproaches []*models.UserApproach) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"updated_at",
		"created_at",
	})
	for _, userApproach := range userApproaches {
		rows.AddRow(
			userApproach.ID,
			userApproach.UpdatedAt,
			userApproach.CreatedAt,
		)
	}
	return rows
}
