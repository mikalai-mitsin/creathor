package postgres

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/interfaces/postgres"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/log"
	"github.com/jmoiron/sqlx"
)

func TestNewPlanRepository(t *testing.T) {
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
		want  repositories.PlanRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &PlanRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPlanRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.plans"
	plan := mock_models.NewPlan(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Plan
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
				mock.ExpectQuery(query).
					WithArgs(
						plan.UpdatedAt,
						plan.CreatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(plan.ID, plan.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						plan.UpdatedAt,
						plan.CreatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PlanRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT plans.id, plans.updated_at, plans.created_at, plans.name, plans.repeat, plans.equipment_id FROM public.plans WHERE id = \\$1 LIMIT 1"
	plan := mock_models.NewPlan(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
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
		want    *models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newPlanRows(t, []*models.Plan{plan})
				mock.ExpectQuery(query).WithArgs(plan.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  plan.ID,
			},
			want:    plan,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(plan.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(plan.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound().WithParam("plan_id", string(plan.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PlanRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanRepository_List(t *testing.T) {
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
	var listPlans []*models.Plan
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listPlans = append(listPlans, mock_models.NewPlan(t))
	}
	filter := mock_models.NewPlanFilter(t)
	query := "SELECT plans.id, plans.updated_at, plans.created_at, plans.name, plans.repeat, plans.equipment_id FROM public.plans"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.PlanFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newPlanRows(t, listPlans))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listPlans,
			wantErr: nil,
		},
		{
			name: "nil page size",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newPlanRows(t, listPlans))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: &models.PlanFilter{},
			},
			want:    listPlans,
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
			r := &PlanRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	plan := mock_models.NewPlan(t)
	query := `UPDATE public.plans`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Plan
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
				mock.ExpectExec(query).
					WithArgs(
						plan.UpdatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
						plan.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						plan.UpdatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
						plan.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: errs.NewEntityNotFound().WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						plan.UpdatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
						plan.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						plan.UpdatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
						plan.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						plan.UpdatedAt,
						plan.Name,
						plan.Repeat,
						plan.EquipmentID,
						plan.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: plan,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PlanRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	plan := mock_models.NewPlan(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
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
				logger:   logger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.plans WHERE id = \\$1").
					WithArgs(plan.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.plans WHERE id = \\$1").
					WithArgs(plan.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			wantErr: errs.NewEntityNotFound().WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.plans WHERE id = \\$1").
					WithArgs(plan.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.plans WHERE id = \\$1").
					WithArgs(plan.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  plan.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("plan_id", string(plan.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PlanRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("PlanRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := `SELECT count\(id\) FROM public.plans`
	ctx := context.Background()
	filter := mock_models.NewPlanFilter(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.PlanFilter
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
			r := &PlanRepository{
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

func newPlanRows(t *testing.T, listPlans []*models.Plan) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"repeat",
		"equipment_id",
		"updated_at",
		"created_at",
	})
	for _, plan := range listPlans {
		rows.AddRow(
			plan.ID,
			plan.Name,
			plan.Repeat,
			plan.EquipmentID,
			plan.UpdatedAt,
			plan.CreatedAt,
		)
	}
	return rows
}
