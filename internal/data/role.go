package data

import (
	"context"
	"rbac/ent"
	"rbac/ent/role"
	"time"
)

type UpdateRoleDto struct {
	Name        string
	Description string
}

type CreateRoleDto struct {
	Name        string
	Description string
	TeamId      int64
}

// RoleRepo
type RoleRepo interface {
	CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error)
	UpdateRole(ctx context.Context, roleId int64, roleDto UpdateRoleDto) (*ent.Role, error)
	DeleteRole(ctx context.Context, roleId int64) error
	GetRoleById(ctx context.Context, roleId int64) (*ent.Role, error)
	GetRoleByIds(ctx context.Context, ids []int64) ([]*ent.Role, error)
	GetRolesList(ctx context.Context, teamId int64, name string) ([]*ent.Role, error)
	AddPermissionToRole(ctx context.Context, roleId int64, permissionId string, fields []string) (*ent.Permission, error)
	RemovePermissionFromRole(ctx context.Context, roleId int64, permissionId string) error
	ListRolePermissions(ctx context.Context, roleId int64) ([]*ent.Permission, error)
}

type roleRepo struct {
	db *ent.Client
}

func validateFields(sourceFields []string, targetFields []string) bool {
	if sourceFields == nil {
		return true
	}
	if targetFields == nil {
		return true
	}
	elementMap := make(map[string]bool)

	// Заполняем карту элементами исходного массива
	for _, element := range sourceFields {
		elementMap[element] = true
	}

	// Проверяем каждый элемент второго массива
	for _, element := range targetFields {
		if _, exists := elementMap[element]; !exists {
			return false
		}
	}

	return true
}

func (r *roleRepo) CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error) {
	return r.db.Role.Create().
		SetName(roleDto.Name).
		SetDescription(roleDto.Description).
		SetTeamID(roleDto.TeamId).Save(ctx)
}

func (r *roleRepo) UpdateRole(ctx context.Context, roleId int64, roleDto UpdateRoleDto) (*ent.Role, error) {
	role, err := r.db.Role.Get(ctx, roleId)
	if err != nil {
		return nil, err
	}
	query := role.Update()
	if roleDto.Name != "" {
		query.SetName(roleDto.Name)
	}
	if roleDto.Description != "" {
		query.SetName(roleDto.Description)
	}
	return query.Save(ctx)
}

func (r *roleRepo) DeleteRole(ctx context.Context, roleId int64) error {
	role, err := r.db.Role.Get(ctx, roleId)
	if err != nil {
		return err
	}
	_, err = role.Update().SetDeletedAt(time.Now()).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *roleRepo) GetRoleById(ctx context.Context, roleId int64) (*ent.Role, error) {
	return r.db.Role.Get(ctx, roleId)
}

func (r *roleRepo) GetRoleByIds(ctx context.Context, ids []int64) ([]*ent.Role, error) {
	return r.db.Role.Query().Where(role.IDIn(ids...)).All(ctx)
}

func (r *roleRepo) GetRolesList(ctx context.Context, teamId int64, name string) ([]*ent.Role, error) {
	query := r.db.Role.Query()
	if name != "" {
		query = query.Where(role.NameContains(name))
	}
	if teamId != 0 {
		query = query.Where(role.TeamID(teamId))
	}
	query = query.Where(role.DeletedAtIsNil())
	return query.All(ctx)
}

func (r *roleRepo) AddPermissionToRole(ctx context.Context, roleId int64, permissionId string, fields []string) (*ent.Permission, error) {
	// check if role exists
	role, err := r.db.Role.Get(ctx, roleId)
	if err != nil {
		return nil, err
	}
	// check if permission exists
	permission, err := r.db.Permission.Get(ctx, permissionId)
	if err != nil {
		return nil, err
	}

	isValid := validateFields(permission.Fields, fields)
	if !isValid {
		panic("Invalid fields")
	}

	_, err = r.db.RolePermission.Create().
		SetFields(fields).
		SetRole(role).
		SetPermission(permission).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *roleRepo) RemovePermissionFromRole(ctx context.Context, roleId int64, permissionId string) error {
	return r.db.RolePermission.DeleteOne(ent.RolePermission{
		ID:           0,
		RoleID:       roleId,
		PermissionID: permissionId,
	}).Exec(ctx)
}

func (r *roleRepo) ListRolePermissions(ctx context.Context, roleId int64) ([]*ent.Permission, error) {
	query := r.db.Role.Query().Where(role.ID(roleId)).QueryPermissions()
	return query.All(ctx)
}

// NewRoleRepo .
func NewRoleRepo(d *Data) RoleRepo {
	return &roleRepo{
		db: d.db,
	}
}
