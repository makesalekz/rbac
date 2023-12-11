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
	Allow       []string
	Deny        []string
}

type CreateRolePermissionDto struct {
	Deny   bool
	Fields []string
}

type FilterRolePermissions struct {
	TenantId    int64
	RolesIds    []int64
	Permissions []string
	DeniedOnly  bool
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
	ListRolesPermissions(ctx context.Context, filter FilterRolePermissions) ([]*ent.RolePermission, error)
}

type roleRepo struct {
	db *ent.Client
}

func (r *roleRepo) CreateRole(ctx context.Context, roleDto CreateRoleDto) (*ent.Role, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := tx.Role.Create().
		SetName(roleDto.Name).
		SetDescription(roleDto.Description).
		SetTenantID(roleDto.TenantId).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	var rpCreate []*ent.RolePermissionCreate

	if len(roleDto.Allow) > 0 {
		for _, rp := range roleDto.Allow {
			query := tx.RolePermission.Create().
				SetTenantID(roleDto.TenantId).
				SetRole(role).
				SetPermissionID(rp).
				SetFields([]string{}).
				SetDeny(false)

			rpCreate = append(rpCreate, query)
		}
	}

	if len(roleDto.Deny) > 0 {
		for _, rp := range roleDto.Deny {
			query := tx.RolePermission.Create().
				SetTenantID(roleDto.TenantId).
				SetRole(role).
				SetPermissionID(rp).
				SetFields([]string{}).
				SetDeny(true)

			rpCreate = append(rpCreate, query)
		}
	}

	if len(rpCreate) > 0 {
		permissions, err := tx.RolePermission.CreateBulk(rpCreate...).Save(ctx)
		if err != nil {
			return nil, err
		}

		role.Edges.Permissions = permissions
	}

	tx.Commit()

	return role, nil
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
		Where(role.ID(roleId), role.TenantIDIn(tenantId, 0)).
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
	return r.db.RolePermission.Create().
		SetTenantID(role.TenantID).
		SetRole(role).
		SetPermission(permission).
		SetFields(dto.Fields).
		SetDeny(dto.Deny).
		Save(ctx)
}

func (r *roleRepo) RemovePermissionFromRole(ctx context.Context, role *ent.Role, permission *ent.Permission) error {
	_, err := r.db.RolePermission.Delete().
		Where(
			rolepermission.TenantID(role.TenantID),
			rolepermission.RoleID(role.ID),
			rolepermission.PermissionID(permission.ID),
		).
		Exec(ctx)

	return err
}

func (r *roleRepo) ListRolePermissions(ctx context.Context, role *ent.Role) ([]*ent.RolePermission, error) {
	return r.db.RolePermission.Query().
		Where(
			rolepermission.TenantID(role.TenantID),
			rolepermission.RoleID(role.ID),
		).
		All(ctx)
}

func (r *roleRepo) ListRolesPermissions(ctx context.Context, filter FilterRolePermissions) ([]*ent.RolePermission, error) {
	query := r.db.RolePermission.
		Query().
		Where(
			rolepermission.RoleIDIn(filter.RolesIds...),
			rolepermission.TenantIDIn(filter.TenantId, 0),
		)

	if len(filter.Permissions) != 0 {
		query.Where(rolepermission.PermissionIDIn(filter.Permissions...))
	}

	if filter.DeniedOnly {
		query.Where(rolepermission.Deny(true))
	}

	return query.All(ctx)
}

// NewRoleRepo .
func NewRoleRepo(d *Data) RoleRepo {
	return &roleRepo{
		db: d.db,
	}
}
