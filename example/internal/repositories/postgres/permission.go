package postgres

import (
	"context"
	"fmt"
	"reflect"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type objectPermissionChecker func(model any, user *models.User) error

var hasObjectPermission = map[models.PermissionID][]objectPermissionChecker{
	models.PermissionIDUserCreate: {objectAnybody},
	models.PermissionIDUserList:   {objectAnybody},
	models.PermissionIDUserDetail: {objectOwner},
	models.PermissionIDUserUpdate: {objectOwner},
	models.PermissionIDUserDelete: {
		objectOwner,
	}, models.PermissionIDSessionList: {objectAnybody}, models.PermissionIDSessionDetail: {objectAnybody}, models.PermissionIDSessionCreate: {objectAnybody}, models.PermissionIDSessionUpdate: {objectAnybody}, models.PermissionIDSessionDelete: {objectAnybody}, models.PermissionIDEquipmentList: {objectAnybody}, models.PermissionIDEquipmentDetail: {objectAnybody}, models.PermissionIDEquipmentCreate: {objectAnybody}, models.PermissionIDEquipmentUpdate: {objectAnybody}, models.PermissionIDEquipmentDelete: {objectAnybody}, models.PermissionIDPlanList: {objectAnybody}, models.PermissionIDPlanDetail: {objectAnybody}, models.PermissionIDPlanCreate: {objectAnybody}, models.PermissionIDPlanUpdate: {objectAnybody}, models.PermissionIDPlanDelete: {objectAnybody}, models.PermissionIDDayList: {objectAnybody}, models.PermissionIDDayDetail: {objectAnybody}, models.PermissionIDDayCreate: {objectAnybody}, models.PermissionIDDayUpdate: {objectAnybody}, models.PermissionIDDayDelete: {objectAnybody}, models.PermissionIDArchList: {objectAnybody}, models.PermissionIDArchDetail: {objectAnybody}, models.PermissionIDArchCreate: {objectAnybody}, models.PermissionIDArchUpdate: {objectAnybody}, models.PermissionIDArchDelete: {objectAnybody},
}

type PermissionRepository struct {
	database *sqlx.DB
}

func NewPermissionRepository(database *sqlx.DB) repositories.PermissionRepository {
	return &PermissionRepository{
		database: database,
	}
}

func (r *PermissionRepository) HasPermission(
	ctx context.Context,
	permissionID models.PermissionID,
	user *models.User,
) error {
	permission := &models.Permission{}
	q := sq.Select("permissions.id", "permissions.name").
		From("public.permissions").
		InnerJoin("group_permissions ON permissions.id = group_permissions.permission_id").
		Where(sq.Eq{"group_permissions.group_id": user.GroupID, "permissions.id": permissionID})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, permission, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(user.ID))
		e.AddParam("permission_id", fmt.Sprint(permissionID))
		return e
	}
	return nil
}

func (r *PermissionRepository) HasObjectPermission(
	_ context.Context,
	permission models.PermissionID,
	user *models.User,
	obj any,
) error {
	checkers := hasObjectPermission[permission]
	for _, checker := range checkers {
		if err := checker(obj, user); err == nil {
			return nil
		}
	}
	return errs.NewPermissionDeniedError()
}

func objectOwner(model any, user *models.User) error {
	valueOf := reflect.ValueOf(model)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return errs.NewPermissionDeniedError()
	}
	modelUserID := valueOf.FieldByName("UserID")
	if modelUserID.Kind() == reflect.Pointer {
		modelUserID = modelUserID.Elem()
	}
	modelID := valueOf.FieldByName("ID")
	if modelID.Kind() == reflect.Pointer {
		modelID = modelID.Elem()
	}
	if modelID.String() == string(user.ID) || modelUserID.String() == string(user.ID) {
		return nil
	}
	return errs.NewPermissionDeniedError()
}

func objectOwnerOrAll(model any, user *models.User) error {
	if model == nil {
		return nil
	}
	valueOf := reflect.ValueOf(model)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return errs.NewPermissionDeniedError()
	}
	modelUserID := valueOf.FieldByName("UserID")
	if modelUserID.Kind() == reflect.Pointer {
		if modelUserID.IsNil() {
			return nil
		}
		modelUserID = modelUserID.Elem()
	}
	modelID := valueOf.FieldByName("ID")
	if modelID.Kind() == reflect.Pointer {
		modelID = modelID.Elem()
	}
	if modelID.String() == string(user.ID) || modelUserID.String() == string(user.ID) {
		return nil
	}
	return errs.NewPermissionDeniedError()
}

func objectNobody(_ any, _ *models.User) error {
	return errs.NewPermissionDeniedError()
}

func objectAnybody(_ any, _ *models.User) error {
	return nil
}
