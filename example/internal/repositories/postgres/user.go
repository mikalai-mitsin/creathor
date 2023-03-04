package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/018bf/example/pkg/postgresql"

	sq "github.com/Masterminds/squirrel"

	"github.com/018bf/example/pkg/log"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserDTO struct {
	ID        string    `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	GroupID   string    `db:"group_id"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}

func NewUserDTOFromModel(user *models.User) *UserDTO {
	dto := &UserDTO{
		ID:        string(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		GroupID:   string(user.GroupID),
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}
	return dto
}

func (dto *UserDTO) ToModel() *models.User {
	model := &models.User{
		ID:        models.UUID(dto.ID),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password,
		Email:     dto.Email,
		GroupID:   models.GroupID(dto.GroupID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
	return model
}

type UserListDTO []*UserDTO

func (list UserListDTO) ToModels() []*models.User {
	listUsers := make([]*models.User, len(list))
	for i := range list {
		listUsers[i] = list[i].ToModel()
	}
	return listUsers
}

type PostgresUserRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPostgresUserRepository(database *sqlx.DB, logger log.Logger) repositories.UserRepository {
	return &PostgresUserRepository{database: database, logger: logger}
}

func (r *PostgresUserRepository) Count(
	ctx context.Context,
	filter *models.UserFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.users")
	if filter.Search != nil {
		q = q.Where(postgresql.Search{
			Lang:   "english",
			Fields: []string{"first_name", "last_name", "email"},
			Query:  *filter.Search,
		})
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result := r.database.QueryRowxContext(ctx, query, args...)
	if err := result.Err(); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewUserDTOFromModel(user)
	q := sq.Insert("public.users").
		Columns(
			"first_name",
			"last_name",
			"password",
			"email",
			"group_id",
			"updated_at",
			"created_at",
		).
		Values(
			dto.FirstName,
			dto.LastName,
			dto.Password,
			dto.Email,
			dto.GroupID,
			dto.UpdatedAt,
			dto.CreatedAt,
		).
		Suffix(`RETURNING id, updated_at, created_at`)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		if errors.Is(e, errs.NewEntityNotFound()) {
			e = errs.NewInvalidFormError()
			e.AddParam("email", "The email field has already been taken.")
		}
		return e
	}
	user.ID = models.UUID(dto.ID)
	user.UpdatedAt = dto.UpdatedAt
	user.CreatedAt = dto.CreatedAt
	return nil
}

func (r *PostgresUserRepository) Get(ctx context.Context, id models.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &UserDTO{}
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &UserDTO{}
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Where(sq.Eq{"email": email}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_email", fmt.Sprint(email))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *PostgresUserRepository) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto UserListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(postgresql.Search{
			Lang:   "english",
			Fields: []string{"first_name", "last_name", "email"},
			Query:  *filter.Search,
		})
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.ToModels(), nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewUserDTOFromModel(user)
	q := sq.Update("public.users").
		Where(sq.Eq{"id": dto.ID}).
		Set("first_name", dto.FirstName).
		Set("last_name", dto.LastName).
		Set("password", dto.Password).
		Set("email", dto.Email).
		Set("group_id", dto.GroupID)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			e = errs.NewInvalidFormError()
			e.AddParam("email", "The email field has already been taken.")
		}
		e.AddParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err)
	}
	if affected == 0 {
		e := errs.NewEntityNotFound()
		e.AddParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.users").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound()
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	return nil
}
