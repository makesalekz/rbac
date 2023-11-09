package data

import (
	"context"
	"rbac/ent"
	"rbac/ent/permission"
)

type UpdatePermissionDto struct {
	Name        string
	Description string
	Fields      []string
}

type CreatePermissionDto struct {
	Id          string
	Name        string
	Description string
	AppId       string
	Fields      []string
}

// PermissionRepo
type PermissionRepo interface {
	CreatePermission(ctx context.Context, createPermissionDto CreatePermissionDto) (*ent.Permission, error)
	UpdatePermission(ctx context.Context, id string, data UpdatePermissionDto) (*ent.Permission, error)
	DeletePermission(ctx context.Context, id string) error
	GetPermissionById(ctx context.Context, id string) (*ent.Permission, error)
	GetPermissionsByIds(ctx context.Context, ids []string) ([]*ent.Permission, error)
	GetPermissions(ctx context.Context, appId string, ids []string) ([]*ent.Permission, error)
}

type permissionRepo struct {
	db *ent.Client
}

func (p *permissionRepo) CreatePermission(ctx context.Context, createPermissionDto CreatePermissionDto) (*ent.Permission, error) {
	return p.db.Permission.Create().
		SetID(createPermissionDto.Id).
		SetName(createPermissionDto.Name).
		SetDescription(createPermissionDto.Description).
		SetAppID(createPermissionDto.AppId).
		SetFields(createPermissionDto.Fields).
		Save(ctx)
}

func (p *permissionRepo) UpdatePermission(ctx context.Context, id string, data UpdatePermissionDto) (*ent.Permission, error) {
	// check permission is exists
	permission, err := p.db.Permission.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	query := permission.Update()
	if data.Name != "" {
		query.SetName(data.Name)
	}
	if data.Description != "" {
		query.SetDescription(data.Description)
	}
	if data.Fields != nil {
		query.SetFields(data.Fields)
	}
	return query.Save(ctx)
}

func (p *permissionRepo) DeletePermission(ctx context.Context, id string) error {
	return p.db.Permission.DeleteOneID(id).Exec(ctx)
}

func (p *permissionRepo) GetPermissionById(ctx context.Context, id string) (*ent.Permission, error) {
	return p.db.Permission.Get(ctx, id)
}

func (p *permissionRepo) GetPermissionsByIds(ctx context.Context, ids []string) ([]*ent.Permission, error) {
	return p.db.Permission.Query().Where(permission.IDIn(ids...)).All(ctx)
}

func (p *permissionRepo) GetPermissions(ctx context.Context, appId string, ids []string) ([]*ent.Permission, error) {
	query := p.db.Permission.Query()
	// check if app exists
	if appId != "" {
		query.Where(permission.AppID(appId))
	}
	if ids != nil {
		query.Where(permission.IDIn(ids...))
	}
	return query.All(ctx)
}

// NewPermissionRepo .
func NewPermissionRepo(d *Data) PermissionRepo {
	return &permissionRepo{
		db: d.db,
	}
}
