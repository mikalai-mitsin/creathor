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

func TestNewMarkRepository(t *testing.T) {
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
		want  repositories.MarkRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &MarkRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewMarkRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarkRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.marks"
	mark := mock_models.NewMark(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Mark
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
						AddRow(mark.ID, mark.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: mark,
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
				card: mark,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &MarkRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT marks.id, marks.name, marks.title, marks.weight, marks.updated_at, marks.created_at FROM public.marks WHERE id = \\$1 LIMIT 1"
	mark := mock_models.NewMark(t)
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
		want    *models.Mark
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newMarkRows(t, []*models.Mark{mark})
				mock.ExpectQuery(query).WithArgs(mark.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  mark.ID,
			},
			want:    mark,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(mark.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("mark_id", mark.ID),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(mark.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound().WithParam("mark_id", mark.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &MarkRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkRepository_List(t *testing.T) {
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
	var marks []*models.Mark
	for i := 0; i < faker.RandomInt(1, 20); i++ {
		marks = append(marks, mock_models.NewMark(t))
	}
	filter := mock_models.NewMarkFilter(t)
	query := "SELECT marks.id, marks.name, marks.title, marks.weight, marks.updated_at, marks.created_at FROM public.marks"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.MarkFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Mark
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newMarkRows(t, marks))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    marks,
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
			r := &MarkRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mark := mock_models.NewMark(t)
	query := `UPDATE public.marks`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Mark
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
				card: mark,
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
				card: mark,
			},
			wantErr: errs.NewEntityNotFound().WithParam("mark_id", mark.ID),
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
				card: mark,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("mark_id", mark.ID),
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
				card: mark,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("mark_id", mark.ID),
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
				card: mark,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("mark_id", mark.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &MarkRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mark := mock_models.NewMark(t)
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
				mock.ExpectExec("DELETE FROM public.marks WHERE id = \\$1").
					WithArgs(mark.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.marks WHERE id = \\$1").
					WithArgs(mark.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			wantErr: errs.NewEntityNotFound().WithParam("mark_id", mark.ID),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.marks WHERE id = \\$1").
					WithArgs(mark.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("mark_id", mark.ID),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.marks WHERE id = \\$1").
					WithArgs(mark.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  mark.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).WithParam("mark_id", mark.ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &MarkRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := `SELECT count\(id\) FROM public.marks`
	ctx := context.Background()
	filter := mock_models.NewMarkFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.MarkFilter
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
			r := &MarkRepository{
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

func newMarkRows(t *testing.T, marks []*models.Mark) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"title",
		"weight",
		"updated_at",
		"created_at",
	})
	for _, mark := range marks {
		rows.AddRow(
			mark.ID,
			mark.Name,
			mark.Title,
			mark.Weight,
			mark.UpdatedAt,
			mark.CreatedAt,
		)
	}
	return rows
}
