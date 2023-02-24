package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/lib/pq"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/interfaces/postgres"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	"github.com/jmoiron/sqlx"
)

func TestNewPostgresUserRepository(t *testing.T) {
	db, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()

	type args struct {
		database *sqlx.DB
		logger   log.Logger
	}
	tests := []struct {
		name string
		args args
		want repositories.UserRepository
	}{
		{
			name: "ok",
			args: args{
				database: db,
			},
			want: &PostgresUserRepository{
				database: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostgresUserRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPostgresUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "INSERT INTO public.users"
	user := mock_models.NewUser(t)
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *models.User
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
				mock.ExpectQuery(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.UpdatedAt,
					user.CreatedAt,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
					AddRow(user.ID, user.CreatedAt))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.UpdatedAt,
					user.CreatedAt,
				).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error": "test error",
				},
			},
		},
		{
			name: "duplicate error",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.UpdatedAt,
					user.CreatedAt,
				).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    3,
				Message: "The form sent is not valid, please correct the errors below.",
				Params: map[string]string{
					"email": "The email field has already been taken.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			if err := r.Create(tt.args.ctx, tt.args.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateByTips() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	user := mock_models.NewUser(t)
	type fields struct {
		database *sqlx.DB
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
			fields: fields{
				database: db,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: &errs.Error{
				Code:    5,
				Message: "Entity not found.",
				Params:  map[string]string{"user_id": string(user.ID)},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = \\$1").
					WithArgs(user.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error":   "test error",
					"user_id": string(user.ID),
				},
			},
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error":   "test error",
					"user_id": string(user.ID),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT id, first_name, last_name, password, email, group_id, created_at, updated_at FROM public.users WHERE id = \\$1 LIMIT 1"
	user := mock_models.NewUser(t)
	type fields struct {
		database *sqlx.DB
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
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := NewUserRows(t, []*models.User{user})
				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    5,
				Message: "Entity not found.",
				Params: map[string]string{
					"user_id": string(user.ID),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
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

func TestPostgresRepository_GetByEmail(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT id, first_name, last_name, password, email, group_id, created_at, updated_at FROM public.users WHERE email = \\$1 LIMIT 1"
	user := mock_models.NewUser(t)
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := NewUserRows(t, []*models.User{user})
				mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:   context.Background(),
				email: user.Email,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:   context.Background(),
				email: user.Email,
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    5,
				Message: "Entity not found.",
				Params: map[string]string{
					"user_email": fmt.Sprint(user.Email),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			got, err := r.GetByEmail(tt.args.ctx, tt.args.email)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByPhone() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_List(t *testing.T) {
	database, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer database.Close()
	pageSize := uint64(100)
	pageNumber := uint64(3)
	var users []*models.User
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		users = append(users, mock_models.NewUser(t))
	}
	query := "SELECT id, first_name, last_name, password, email, group_id, created_at, updated_at FROM public.users"
	filter := mock_models.NewUserFilter(t)
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		filter *models.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(NewUserRows(t, users))
			},
			fields: fields{
				database: database,
			},
			args: args{
				ctx:    context.Background(),
				filter: filter,
			},
			want:    users,
			wantErr: nil,
		},
		{
			name: "nil page",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(NewUserRows(t, users))
			},
			fields: fields{
				database: database,
			},
			args: args{
				ctx:    context.Background(),
				filter: &models.UserFilter{},
			},
			want:    users,
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: database,
			},
			args: args{
				ctx: context.Background(),
				filter: &models.UserFilter{
					PageSize:   &pageSize,
					PageNumber: &pageNumber,
					Search:     utils.Pointer("asd"),
					OrderBy:    []string{"id", "created_at"},
				},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
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

func TestPostgresRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	user := mock_models.NewUser(t)
	query := `UPDATE public.users`
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *models.User
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
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    5,
				Message: "Entity not found.",
				Params: map[string]string{
					"user_id": string(user.ID),
				},
			},
		},
		{
			name: "duplicate error",
			setup: func() {
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnError(&pq.Error{
					Code: "23505",
				})
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    3,
				Message: "The form sent is not valid, please correct the errors below.",
				Params: map[string]string{
					"user_id": string(user.ID),
					"email":   "The email field has already been taken.",
				},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"user_id": string(user.ID),
					"error":   "test error",
				},
			},
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error":   "test error",
					"user_id": string(user.ID),
				},
			},
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).WithArgs(
					user.FirstName,
					user.LastName,
					user.Password,
					user.Email,
					user.GroupID,
					user.ID,
				).WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error": "test error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
				database: tt.fields.database,
			}
			if err := r.Update(tt.args.ctx, tt.args.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresUserRepository_Count(t *testing.T) {
	database, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer database.Close()
	query := `SELECT count\(id\) FROM public.users`
	filter := mock_models.NewUserFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserFilter
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
				mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			fields: fields{
				database: database,
			},
			args: args{
				ctx:    context.Background(),
				filter: filter,
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "bad return type",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("one"))
			},
			fields: fields{
				database: database,
			},
			args: args{
				ctx: context.Background(),
				filter: &models.UserFilter{
					PageSize:   utils.Pointer(uint64(10)),
					PageNumber: utils.Pointer(uint64(2)),
					Search:     utils.Pointer("e"),
					OrderBy:    []string{"created_at"},
				},
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
				database: database,
			},
			args: args{
				ctx: context.Background(),
				filter: &models.UserFilter{
					PageSize:   nil,
					PageNumber: nil,
				},
			},
			want: 0,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error": "test error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := PostgresUserRepository{
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

func NewUserRows(t *testing.T, users []*models.User) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID,
			user.FirstName,
			user.LastName,
			user.Password,
			user.Email,
			user.GroupID,
			user.CreatedAt,
			user.UpdatedAt,
		)
	}
	return rows
}
