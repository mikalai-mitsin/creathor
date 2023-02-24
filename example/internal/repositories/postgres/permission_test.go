package postgres

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/interfaces/postgres"
	"github.com/018bf/example/pkg/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestPermissionRepository_objectAnybody(t *testing.T) {
	type fields struct {
	}
	type args struct {
		in0 any
		in1 *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
			},
			fields:  fields{},
			args:    args{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if err := objectAnybody(tt.args.in0, tt.args.in1); !errors.Is(err, tt.wantErr) {
				t.Errorf("objectAnybody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionRepository_objectNobody(t *testing.T) {
	type fields struct {
	}
	type args struct {
		in0 any
		in1 *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
			},
			fields:  fields{},
			args:    args{},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if err := objectNobody(tt.args.in0, tt.args.in1); !errors.Is(err, tt.wantErr) {
				t.Errorf("objectNobody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionRepository_objectOwner(t *testing.T) {
	db, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		model any
		user  *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-241"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name: "no struct",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				model: "&models.Tips{}",
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "self",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-241"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name: "self id pointer",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				model: struct {
					ID *string
				}{
					ID: utils.Pointer("asd-241"),
				},
				user: &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name: "permission denied",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-2412"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := objectOwner(tt.args.model, tt.args.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("objectOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPermissionRepository(t *testing.T) {
	db, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	type args struct {
		database *sqlx.DB
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  repositories.PermissionRepository
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				database: db,
			},
			want: &PermissionRepository{database: db},
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPermissionRepository(tt.args.database); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPermissionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionRepository_HasPermission(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT permissions.id, permissions.name FROM public.permissions INNER JOIN group_permissions ON permissions.id = group_permissions.permission_id WHERE group_permissions.group_id = \\$1 AND permissions.id = \\$2"
	user := mock_models.NewUser(t)
	user.GroupID = "user"
	permission := &models.Permission{}
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		ctx          context.Context
		permissionID models.PermissionID
		user         *models.User
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
				rows := NewPermissionRows(t, []*models.Permission{permission})
				mock.ExpectQuery(query).
					WithArgs(user.GroupID, models.PermissionIDUserCreate).
					WillReturnRows(rows)
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:          context.Background(),
				permissionID: models.PermissionIDUserCreate,
				user:         user,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(user.GroupID, models.PermissionIDUserList).
					WillReturnError(errors.New("error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:          context.Background(),
				permissionID: models.PermissionIDUserList,
				user:         user,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"error":         "error",
					"permission_id": "user_list",
					"user_id":       string(user.ID),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PermissionRepository{
				database: tt.fields.database,
			}
			if err := r.HasPermission(tt.args.ctx, tt.args.permissionID, tt.args.user); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("HasPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionRepository_HasObjectPermission1(t *testing.T) {
	db, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	user := mock_models.NewUser(t)
	user.GroupID = "user"
	article := mock_models.NewUser(t)
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		in0        context.Context
		permission models.PermissionID
		user       *models.User
		obj        any
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
			},
			fields: fields{
				database: db,
			},
			args: args{
				in0:        nil,
				permission: models.PermissionIDUserCreate,
				user:       user,
				obj:        nil,
			},
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
			},
			fields: fields{
				database: db,
			},
			args: args{
				in0:        nil,
				permission: models.PermissionIDUserUpdate,
				user:       user,
				obj:        article,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PermissionRepository{
				database: tt.fields.database,
			}
			if err := r.HasObjectPermission(tt.args.in0, tt.args.permission, tt.args.user, tt.args.obj); !errors.Is(
				err,
				tt.wantErr,
			) {
				t.Errorf("HasObjectPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionRepository_objectOwnerOrAll(t *testing.T) {
	db, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	type fields struct {
		database *sqlx.DB
	}
	type args struct {
		model any
		user  *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:  "ok",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-241"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name:  "no struct",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: "&models.Tips{}",
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name:  "self",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-241"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name:  "self id pointer",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: struct {
					ID *string
				}{
					ID: utils.Pointer("asd-241"),
				},
				user: &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
		{
			name:  "permission denied",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: &models.User{ID: "asd-2412"},
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name:  "model is nil",
			setup: func() {},
			fields: fields{
				database: db,
			},
			args: args{
				model: nil,
				user:  &models.User{ID: "asd-241"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if err := objectOwnerOrAll(tt.args.model, tt.args.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("objectOwnerOrAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func NewPermissionRows(t *testing.T, permission []*models.Permission) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
	})
	for _, perm := range permission {
		rows.AddRow(
			perm.ID,
			perm.Name,
		)
	}
	return rows
}
