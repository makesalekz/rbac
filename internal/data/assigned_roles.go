package data

import (
	"context"

	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/predicate"
	"gitlab.calendaria.team/services/rbac/ent/resourceaccess"

	_ "github.com/lib/pq"
)

// AssignedRolesRepo
type AssignedRolesRepo interface {
	AssignRoles(ctx context.Context, tenantId int64, dto []AssignRoleDto) error
	UnassignRole(ctx context.Context, assignedRole *ent.ResourceAccess) error
	GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.ResourceAccess, error)
	ListAssignedRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error)
	CheckRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error)
}

type assignedRolesRepo struct {
	db *ent.Client
}

func (t *assignedRolesRepo) AssignRoles(ctx context.Context, tenantId int64, dtos []AssignRoleDto) error {
	assignQueries := make([]*ent.ResourceAccessCreate, len(dtos))
	for i, dto := range dtos {
		assignQueries[i] = t.db.ResourceAccess.Create().
			SetTenantID(tenantId).
			SetIdentityID(dto.IdentityID).
			SetRoleID(dto.RoleID)

		if dto.Resource != nil {
			assignQueries[i].SetResourceID(dto.Resource.Id).SetResourceType(dto.Resource.Type)
		}
	}

	return t.db.ResourceAccess.CreateBulk(assignQueries...).Exec(ctx)
}

func (t *assignedRolesRepo) GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.ResourceAccess, error) {
	return t.db.ResourceAccess.Query().
		Where(
			resourceaccess.TenantID(tenantId),
			resourceaccess.ID(assignId),
		).
		Only(ctx)
}

func (t *assignedRolesRepo) UnassignRole(ctx context.Context, assignedRole *ent.ResourceAccess) error {
	return t.db.ResourceAccess.DeleteOne(assignedRole).Exec(ctx)
}

func (t *assignedRolesRepo) ListAssignedRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error) {
	query := t.db.ResourceAccess.Query().
		Where(resourceaccess.TenantID(dto.TenantID)).
		WithRole()

	if len(dto.IdentityIDs) > 0 {
		identityIDs := append(dto.IdentityIDs, "") // all provided identities + "all" identity
		query.Where(resourceaccess.IdentityIDIn(identityIDs...))
	}

	if len(dto.ResourceFilter) > 0 {
		query.Where(
			resourceaccess.Or(
				resourceaccess.ResourceTypeIn(dto.ResourceFilter...),
				resourceaccess.ResourceIDIsNil(),
			),
		)
	} else if len(dto.Resources) > 0 {
		// assigned only on provided resource
		predicates := make([]predicate.ResourceAccess, len(dto.Resources))
		for i, resource := range dto.Resources {
			predicates[i] = resourceaccess.And(
				resourceaccess.ResourceType(resource.Type),
				resourceaccess.ResourceID(resource.Id),
			)
		}

		query.Where(resourceaccess.Or(predicates...))
	} else {
		// not assigned on any resource (all resources on tenant)
		query.Where(resourceaccess.ResourceIDIsNil())
	}

	return query.All(ctx)
}

func (t *assignedRolesRepo) CheckRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error) {
	query := t.db.ResourceAccess.Query().
		Where(resourceaccess.TenantID(dto.TenantID)).
		WithRole()

	if len(dto.IdentityIDs) > 0 {
		identityIDs := append(dto.IdentityIDs, "")
		query.Where(resourceaccess.IdentityIDIn(identityIDs...))
	}

	predicates := []predicate.ResourceAccess{
		resourceaccess.ResourceIDIsNil(),
	}

	if len(dto.Resources) > 0 {
		for _, resource := range dto.Resources {
			predicates = append(predicates,
				resourceaccess.And(
					resourceaccess.ResourceType(resource.Type),
					resourceaccess.ResourceID(resource.Id),
				),
			)
		}
	}

	return query.Where(resourceaccess.Or(predicates...)).All(ctx)
}

// NewAssignedRolesRepo .
func NewAssignedRolesRepo(d *Data) AssignedRolesRepo {
	return &assignedRolesRepo{
		db: d.db,
	}
}
