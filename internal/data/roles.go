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
	Allow       []string
	Deny        []string
}

type CreateRoleDto struct {
	Name        string
	Description string
	TenantId    int64
	IsSystem    bool
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
	SetRolePermission(ctx context.Context, role *ent.Role, permission *ent.Permission, dto CreateRolePermissionDto) error
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
		SetIsSystem(roleDto.IsSystem).
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
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// update role data
	shouldUpdate := false
	query := tx.Role.UpdateOne(role)

	if roleDto.Name != role.Name {
		query.SetName(roleDto.Name)
		shouldUpdate = true
	}

	if roleDto.Description != role.Description {
		query.SetDescription(roleDto.Description)
		shouldUpdate = true
	}

	if shouldUpdate {
		role, err = query.Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// get existing role permissions
	permissions, err := tx.RolePermission.Query().
		Where(
			rolepermission.TenantID(role.TenantID),
			rolepermission.RoleID(role.ID),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var rpCreate []*ent.RolePermissionCreate
	var rpDelete []string

	allow := make(map[string]bool)
	deny := make(map[string]bool)

	for _, pid := range roleDto.Allow {
		allow[pid] = true
	}

	for _, pid := range roleDto.Deny {
		deny[pid] = true
	}

	// fill rpCreate and rpDelete
	for _, rp := range permissions {
		// don't rewrite existing permissions
		if allow[rp.PermissionID] {
			if !rp.Deny {
				delete(allow, rp.PermissionID)
			}
			continue
		}

		if deny[rp.PermissionID] {
			if rp.Deny {
				delete(deny, rp.PermissionID)
			}
			continue
		}

		rpDelete = append(rpDelete, rp.PermissionID)
	}

	for pid := range allow {
		query := tx.RolePermission.Create().
			SetTenantID(role.TenantID).
			SetRole(role).
			SetPermissionID(pid).
			SetFields([]string{}).
			SetDeny(false)

		rpCreate = append(rpCreate, query)
	}

	for pid := range deny {
		query := tx.RolePermission.Create().
			SetTenantID(role.TenantID).
			SetRole(role).
			SetPermissionID(pid).
			SetFields([]string{}).
			SetDeny(true)

		rpCreate = append(rpCreate, query)
	}

	// update role permissions
	if len(rpDelete) > 0 {
		_, err := tx.RolePermission.Delete().
			Where(
				rolepermission.RoleID(role.ID),
				rolepermission.PermissionIDIn(rpDelete...),
			).
			Exec(ctx)

		if err != nil {
			return nil, err
		}
	}

	if len(rpCreate) > 0 {
		err := tx.RolePermission.CreateBulk(rpCreate...).
			OnConflictColumns(rolepermission.FieldRoleID, rolepermission.FieldPermissionID).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	tx.Commit()

	return role, nil
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

func (r *roleRepo) SetRolePermission(ctx context.Context, role *ent.Role, permission *ent.Permission, dto CreateRolePermissionDto) error {
	return r.db.RolePermission.Create().
		SetTenantID(role.TenantID).
		SetRole(role).
		SetPermission(permission).
		SetFields(dto.Fields).
		SetDeny(dto.Deny).
		OnConflictColumns(rolepermission.FieldRoleID, rolepermission.FieldPermissionID).
		UpdateNewValues().
		Exec(ctx)
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
