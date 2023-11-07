package data

import (
	"context"
	"rbac/ent"
	"rbac/ent/role"
	"rbac/ent/rolepermission"
)

type UpdateRoleDto struct {
	Name        string
	Description string
}

type CreateRoleDto struct {
	Name        string
	Description string
	TenantId    int64
}

type CreateRolePermissionDto struct {
	RoleId       int64
	PermissionId string
	TenantId     int64
	Deny         bool
	Fields       []string
}

// RoleRepo
type RoleRepo interface {
	CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error)
	UpdateRole(ctx context.Context, roleID int64, roleDto UpdateRoleDto) (*ent.Role, error)
	DeleteRole(ctx context.Context, roleID int64) error
	GetRoleById(ctx context.Context, roleID, tenantId int64) (*ent.Role, error)
	GetRoleByIds(ctx context.Context, ids []int64) ([]*ent.Role, error)
	GetRolesList(ctx context.Context, teamID int64, name string) ([]*ent.Role, error)
	AddPermissionToRole(ctx context.Context, dto CreateRolePermissionDto) (*ent.RolePermission, error)
	RemovePermissionFromRole(ctx context.Context, roleID, tenantId int64, permissionId string) error
	ListRolePermissions(ctx context.Context, roleId, tenantId int64) ([]*ent.RolePermission, error)
	ListRolesPermissions(ctx context.Context, roleId []int64, tenantId int64, permissions []string) ([]*ent.RolePermission, error)
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
		SetTenantID(roleDto.TenantId).Save(ctx)
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
	return r.db.Role.DeleteOneID(roleId).Exec(ctx)
}

func (r *roleRepo) GetRoleById(ctx context.Context, roleId, tenantId int64) (*ent.Role, error) {
	query := r.db.Role.Query()
	query = query.Where(role.ID(roleId), role.TenantID(tenantId))
	return query.First(ctx)
}

func (r *roleRepo) GetRoleByIds(ctx context.Context, ids []int64) ([]*ent.Role, error) {
	return r.db.Role.Query().Where(role.IDIn(ids...)).All(ctx)
}

func (r *roleRepo) GetRolesList(ctx context.Context, tenantId int64, name string) ([]*ent.Role, error) {
	query := r.db.Role.Query()
	if name != "" {
		query = query.Where(role.NameContains(name))
	}
	if tenantId != 0 {
		query = query.Where(role.TenantID(tenantId))
	}
	return query.All(ctx)
}

func (r *roleRepo) AddPermissionToRole(ctx context.Context, dto CreateRolePermissionDto) (*ent.RolePermission, error) {
	// check if role exists
	role, err := r.db.Role.Get(ctx, dto.RoleId)
	if err != nil {
		return nil, err
	}
	// check if permission exists
	permission, err := r.db.Permission.Get(ctx, dto.PermissionId)
	if err != nil {
		return nil, err
	}

	isValid := validateFields(permission.Fields, dto.Fields)
	if !isValid {
		panic("Invalid fields")
	}

	rolePermissionSave, err := r.db.RolePermission.Create().
		SetFields(dto.Fields).
		SetRole(role).
		SetPermission(permission).
		SetTenantID(dto.TenantId).
		SetDeny(dto.Deny).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return rolePermissionSave, nil
}

func (r *roleRepo) RemovePermissionFromRole(ctx context.Context, roleId, tenantId int64, permissionId string) error {
	return r.db.RolePermission.DeleteOne(&ent.RolePermission{
		RoleID:       roleId,
		PermissionID: permissionId,
		TenantID:     &tenantId,
	}).Exec(ctx)
}

func (r *roleRepo) ListRolePermissions(ctx context.Context, roleId, tenantId int64) ([]*ent.RolePermission, error) {
	query := r.db.RolePermission.
		Query().
		Where(
			rolepermission.HasRoleWith(role.ID(roleId)),
			rolepermission.TenantIDEQ(tenantId),
		)

	return query.All(ctx)
}

func (r *roleRepo) ListRolesPermissions(ctx context.Context, roleIds []int64, tenantId int64, permissions []string) ([]*ent.RolePermission, error) {
	query := r.db.RolePermission.
		Query().
		Where(
			rolepermission.HasRoleWith(role.IDIn(roleIds...)),
			rolepermission.TenantIDEQ(tenantId),
			rolepermission.TenantID(0),
		)
	if len(permissions) != 0 {
		query = query.Where(rolepermission.PermissionIDIn(permissions...))
	}
	return query.All(ctx)
}

// NewRoleRepo .
func NewRoleRepo(d *Data) RoleRepo {
	return &roleRepo{
		db: d.db,
	}
}
