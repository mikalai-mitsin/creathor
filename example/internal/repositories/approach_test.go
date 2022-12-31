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

func TestNewPostgresApproachRepository(t *testing.T) {
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
		want  repositories.ApproachRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &PostgresApproachRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPostgresApproachRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostgresApproachRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresApproachRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.approachs"
	approach := mock_models.NewApproach(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Approach
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
						AddRow(approach.ID, approach.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: approach,
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
				card: approach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostgresApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("PostgresApproachRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresApproachRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT approachs.id, approachs.updated_at, approachs.created_at FROM public.approachs WHERE id = \\$1 LIMIT 1"
	approach := mock_models.NewApproach(t)
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
		want    *models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newApproachRows(t, []*models.Approach{approach})
				mock.ExpectQuery(query).WithArgs(approach.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  approach.ID,
			},
			want:    approach,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(approach.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("approach_id", approach.ID),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(approach.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound().WithParam("approach_id", approach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostgresApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostgresApproachRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresApproachRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresApproachRepository_List(t *testing.T) {
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
	var approachs []*models.Approach
	for i := 0; i < faker.RandomInt(1, 20); i++ {
		approachs = append(approachs, mock_models.NewApproach(t))
	}
	filter := mock_models.NewApproachFilter(t)
	query := "SELECT approachs.id, approachs.updated_at, approachs.created_at FROM public.approachs"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ApproachFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Approach
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newApproachRows(t, approachs))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    approachs,
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
			r := &PostgresApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostgresApproachRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresApproachRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresApproachRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	approach := mock_models.NewApproach(t)
	query := `UPDATE public.approachs`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Approach
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
				card: approach,
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
				card: approach,
			},
			wantErr: errs.NewEntityNotFound().WithParam("approach_id", approach.ID),
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
				card: approach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("approach_id", approach.ID),
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
				card: approach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("approach_id", approach.ID),
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
				card: approach,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("approach_id", approach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostgresApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("PostgresApproachRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresApproachRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	approach := mock_models.NewApproach(t)
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
				mock.ExpectExec("DELETE FROM public.approachs WHERE id = \\$1").
					WithArgs(approach.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.approachs WHERE id = \\$1").
					WithArgs(approach.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			wantErr: errs.NewEntityNotFound().WithParam("approach_id", approach.ID),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.approachs WHERE id = \\$1").
					WithArgs(approach.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("approach_id", approach.ID),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.approachs WHERE id = \\$1").
					WithArgs(approach.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  approach.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("approach_id", approach.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostgresApproachRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("PostgresApproachRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresApproachRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := `SELECT count\(id\) FROM public.approachs`
	ctx := context.Background()
	filter := mock_models.NewApproachFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ApproachFilter
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
			r := &PostgresApproachRepository{
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

func newApproachRows(t *testing.T, approachs []*models.Approach) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"updated_at",
		"created_at",
	})
	for _, approach := range approachs {
		rows.AddRow(
			approach.ID,
			approach.UpdatedAt,
			approach.CreatedAt,
		)
	}
	return rows
}
