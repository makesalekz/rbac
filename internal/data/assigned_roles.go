package data

import (
	"context"

	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/predicate"
	"gitlab.calendaria.team/services/rbac/ent/resourceaccess"

	_ "github.com/lib/pq"
)

// AssignedRolesRepo.
type AssignedRolesRepo interface {
	AssignRoles(ctx context.Context, tenantID int64, dto []AssignRoleDto) error
	UnassignRole(ctx context.Context, assignedRole *ent.ResourceAccess) error
	GetAssignedRoleByID(ctx context.Context, tenantID, assignID int64) (*ent.ResourceAccess, error)
	ListAssignedRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error)
	CheckRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error)
}

type assignedRolesRepo struct {
	db *ent.Client
}

func (t *assignedRolesRepo) AssignRoles(ctx context.Context, tenantID int64, dtos []AssignRoleDto) error {
	assignQueries := make([]*ent.ResourceAccessCreate, len(dtos))
	for i, dto := range dtos {
		assignQueries[i] = t.db.ResourceAccess.Create().
			SetTenantID(tenantID).
			SetIdentityID(dto.IdentityID).
			SetRoleID(dto.RoleID)

		if dto.Resource != nil {
			assignQueries[i].SetResourceID(dto.Resource.GetId()).SetResourceType(dto.Resource.GetType())
		}
	}

	return t.db.ResourceAccess.CreateBulk(assignQueries...).Exec(ctx)
}

func (t *assignedRolesRepo) GetAssignedRoleByID(
	ctx context.Context,
	tenantID, assignID int64,
) (*ent.ResourceAccess, error) {
	return t.db.ResourceAccess.Query().
		Where(
			resourceaccess.TenantID(tenantID),
			resourceaccess.ID(assignID),
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
		dto.IdentityIDs = append(dto.IdentityIDs, "") // all provided identities + "all" identity
		query.Where(resourceaccess.IdentityIDIn(dto.IdentityIDs...))
	}

	switch {
	case len(dto.ResourceFilter) > 0:
		query.Where(
			resourceaccess.Or(
				resourceaccess.ResourceTypeIn(dto.ResourceFilter...),
				resourceaccess.ResourceIDIsNil(),
			),
		)
	case len(dto.Resources) > 0:
		// assigned only on provided resource
		predicates := make([]predicate.ResourceAccess, len(dto.Resources))
		for i, resource := range dto.Resources {
			predicates[i] = resourceaccess.And(
				resourceaccess.ResourceType(resource.GetType()),
				resourceaccess.ResourceID(resource.GetId()),
			)
		}

		query.Where(resourceaccess.Or(predicates...))
	default:
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
		dto.IdentityIDs = append(dto.IdentityIDs, "")
		query.Where(resourceaccess.IdentityIDIn(dto.IdentityIDs...))
	}

	predicates := []predicate.ResourceAccess{
		resourceaccess.ResourceIDIsNil(),
	}

	if len(dto.Resources) > 0 {
		for _, resource := range dto.Resources {
			predicates = append(predicates,
				resourceaccess.And(
					resourceaccess.ResourceType(resource.GetType()),
					resourceaccess.ResourceID(resource.GetId()),
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
