package data

import (
	"context"

	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/role"
	"gitlab.calendaria.team/services/rbac/ent/rolepermission"
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
	Deny   bool
	Fields []string
}

// RoleRepo
type RoleRepo interface {
	CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error)
	UpdateRole(ctx context.Context, role *ent.Role, roleDto UpdateRoleDto) (*ent.Role, error)
	DeleteRole(ctx context.Context, role *ent.Role) error
	GetRoleById(ctx context.Context, tenantId, roleId int64) (*ent.Role, error)
	GetRolesList(ctx context.Context, tenantId int64, search string) ([]*ent.Role, error)
	AddPermissionToRole(ctx context.Context, role *ent.Role, permission *ent.Permission, dto CreateRolePermissionDto) (*ent.RolePermission, error)
	RemovePermissionFromRole(ctx context.Context, role *ent.Role, permission *ent.Permission) error
	ListRolePermissions(ctx context.Context, role *ent.Role) ([]*ent.RolePermission, error)
	ListRolesPermissions(ctx context.Context, roleId []int64, tenantId int64, permissions []string) ([]*ent.RolePermission, error)
}

type roleRepo struct {
	db *ent.Client
}

func (r *roleRepo) CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error) {
	return r.db.Role.Create().
		SetName(roleDto.Name).
		SetDescription(roleDto.Description).
		SetTenantID(roleDto.TenantId).
		Save(ctx)
}

func (r *roleRepo) UpdateRole(ctx context.Context, role *ent.Role, roleDto UpdateRoleDto) (*ent.Role, error) {
	query := r.db.Role.UpdateOne(role)

	if roleDto.Name != "" {
		query.SetName(roleDto.Name)
	}
	if roleDto.Description != "" {
		query.SetDescription(roleDto.Description)
	}

	return query.Save(ctx)
}

func (r *roleRepo) DeleteRole(ctx context.Context, role *ent.Role) error {
	return r.db.Role.DeleteOne(role).Exec(ctx)
}

func (r *roleRepo) GetRoleById(ctx context.Context, tenantId, roleId int64) (*ent.Role, error) {
	return r.db.Role.Query().
		Where(role.ID(roleId), role.TenantID(tenantId)).
		First(ctx)
}

func (r *roleRepo) GetRolesList(ctx context.Context, tenantId int64, search string) ([]*ent.Role, error) {
	query := r.db.Role.Query().Where(role.TenantID(tenantId))

	if search != "" {
		query = query.Where(role.NameContainsFold(search))
	}

	return query.All(ctx)
}

func (r *roleRepo) AddPermissionToRole(ctx context.Context, role *ent.Role, permission *ent.Permission, dto CreateRolePermissionDto) (*ent.RolePermission, error) {
	query := r.db.RolePermission.Create().
		SetRole(role).
		SetPermission(permission).
		SetFields(dto.Fields).
		SetDeny(dto.Deny)

	if role.TenantID != nil {
		query.SetTenantID(*role.TenantID)
	}

	return query.Save(ctx)
}

func (r *roleRepo) RemovePermissionFromRole(ctx context.Context, role *ent.Role, permission *ent.Permission) error {
	query := r.db.RolePermission.Delete().
		Where(
			rolepermission.RoleID(role.ID),
			rolepermission.PermissionID(permission.ID),
		)

	if role.TenantID != nil {
		query.Where(rolepermission.TenantID(*role.TenantID))
	}

	_, err := query.Exec(ctx)

	return err
}

func (r *roleRepo) ListRolePermissions(ctx context.Context, role *ent.Role) ([]*ent.RolePermission, error) {
	query := r.db.RolePermission.Query().
		Where(rolepermission.RoleID(role.ID))

	if role.TenantID != nil {
		query.Where(rolepermission.TenantID(*role.TenantID))
	}

	return query.All(ctx)
}

func (r *roleRepo) ListRolesPermissions(ctx context.Context, roleIds []int64, tenantId int64, permissions []string) ([]*ent.RolePermission, error) {
	query := r.db.RolePermission.
		Query().
		Where(
			rolepermission.RoleIDIn(roleIds...),
			rolepermission.TenantIDIn(tenantId, 0),
		)

	if len(permissions) != 0 {
		query.Where(rolepermission.PermissionIDIn(permissions...))
	}

	return query.All(ctx)
}

// NewRoleRepo .
func NewRoleRepo(d *Data) RoleRepo {
	return &roleRepo{
		db: d.db,
	}
}
