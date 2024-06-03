package data

import (
	"context"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/predicate"
	"gitlab.calendaria.team/services/rbac/ent/resourceaccess"

	_ "github.com/lib/pq"
)

const RESOURCE_TYPE_TEAM = "team"

type AssignRoleDto struct {
	IdentityId string
	RoleId     int64
	TeamId     int64
	Resource   *v1.Resource
}

type ListRolesDto struct {
	TenantId       int64
	IdentityIDs    []string
	TeamsIDs       []int64
	Resources      []*v1.Resource
	ResourceFilter *v1.ResourceFilter
}

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
			SetIdentityID(dto.IdentityId).
			SetRoleID(dto.RoleId)

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
		Where(resourceaccess.TenantID(dto.TenantId)).
		WithRole()

	if len(dto.IdentityIDs) > 0 {
		identityIDs := append(dto.IdentityIDs, "") // all provided identities + "all" identity
		query.Where(resourceaccess.IdentityIDIn(identityIDs...))
	}

	if dto.ResourceFilter != nil && len(dto.ResourceFilter.ResourceTypes) > 0 {
		query.Where(
			resourceaccess.Or(
				resourceaccess.ResourceTypeIn(dto.ResourceFilter.ResourceTypes...),
				resourceaccess.ResourceIDIsNil(),
			),
		)
	} else if len(dto.Resources) > 0 {
		// assigned only on provided resource
		for _, resource := range dto.Resources {
			query.Where(
				resourceaccess.ResourceType(resource.Type),
				resourceaccess.ResourceID(resource.Id),
			)
		}
	} else {
		// not assigned on any resource (all resources on tenant)
		query.Where(resourceaccess.ResourceIDIsNil())
	}

	return query.All(ctx)
}

func (t *assignedRolesRepo) CheckRoles(ctx context.Context, dto ListRolesDto) ([]*ent.ResourceAccess, error) {
	query := t.db.ResourceAccess.Query().
		Where(resourceaccess.TenantID(dto.TenantId)).
		WithRole()

	if len(dto.IdentityIDs) > 0 {
		identityIDs := append(dto.IdentityIDs, "")
		query.Where(resourceaccess.IdentityIDIn(identityIDs...))
	}

	predicates := []predicate.ResourceAccess{
		resourceaccess.ResourceIDIsNil(),
	}

	if len(dto.TeamsIDs) > 0 {
		predicates = append(predicates,
			resourceaccess.And(
				resourceaccess.ResourceType(RESOURCE_TYPE_TEAM),
				resourceaccess.ResourceIDIn(dto.TeamsIDs...),
			),
		)
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
