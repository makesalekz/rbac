package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	u_nats "gitlab.calendaria.team/services/utils/v1/nats"
)

type AssignRoleMessage struct {
	data.AssignRoleDto
	TenantID int64
}

// AssignedRolesUsecase .
type AssignedRolesUsecase struct {
	log      *log.Helper
	repo     data.AssignedRolesRepo
	roleRepo data.RoleRepo
	teamRepo data.TeamsRepo
	qm       *u_nats.QueueManager
}

// NewAssignedRolesUsecase .
func NewAssignedRolesUsecase(
	logger log.Logger,
	repo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamRepo data.TeamsRepo,
	qm *u_nats.QueueManager,
) (*AssignedRolesUsecase, error) {
	return &AssignedRolesUsecase{
		log:      log.NewHelper(log.With(logger, "module", "biz/rbac")),
		repo:     repo,
		roleRepo: roleRepo,
		teamRepo: teamRepo,
		qm:       qm,
	}, nil
}

// AssignRoles assigns multiple roles to identities.
// This function checks if the role exists and if the team exists (if team id is not zero).
//
// Possible errors:
// - rbac.ErrorNotFound: role or team not found
// - rbac.ErrorDatabaseQuery: failed to get role or team
// - rbac.ErrorAlreadyExists: role already assigned
// - rbac.ErrorBadRequest: there is no such teamId.
func (u *AssignedRolesUsecase) AssignRoles(ctx context.Context, tenantID int64, dtos []data.AssignRoleDto) error {
	roleIDs := data.ExtractUnique(dtos, func(e data.AssignRoleDto) (int64, bool) { return e.RoleId, true })
	// Get roles by ids
	roles, err := u.roleRepo.GetRolesByID(ctx, tenantID, roleIDs)
	if err != nil {
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}
	if len(roles) < len(roleIDs) {
		foundIDs := data.ExtractSlice(roles, func(e *ent.Role) (int64, bool) { return e.ID, true })
		diff := data.Diff(roleIDs, foundIDs)
		return v1.ErrorBadRequest("invalid role ids %v", diff)
	}

	teamIDs := data.ExtractUnique(dtos, func(e data.AssignRoleDto) (int64, bool) {
		if e.TeamId != 0 {
			return e.TeamId, true
		}
		return 0, false
	})

	// If there are team ids, get teams by ids, then check if the returned ids are equal to the input ids
	if len(teamIDs) > 0 {
		teams, err := u.teamRepo.GetTeams(ctx, tenantID, teamIDs)
		if err != nil {
			return v1.ErrorDatabaseQuery("get teams failed: %v", err)
		}
		if len(teams) < len(teamIDs) {
			foundIDs := data.ExtractSlice(teams, func(e *ent.Team) (int64, bool) { return e.ID, true })
			diff := data.Diff(teamIDs, foundIDs)
			return v1.ErrorBadRequest("invalid team ids %v", diff)
		}
	}

	// Assign roles
	err = u.repo.AssignRoles(ctx, tenantID, dtos)
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	// TODO: remove check, with mock for QueueManager
	if u.qm != nil {
		queue := u.qm.GetLocal(QueueRoleAssign)

		for _, dto := range dtos {
			queue.Pub(AssignRoleMessage{
				AssignRoleDto: dto,
				TenantID:      tenantID,
			})
		}
	}

	return nil
}

func (u *AssignedRolesUsecase) AssignRole(ctx context.Context, tenantID int64, dto data.AssignRoleDto) error {
	_, err := u.roleRepo.GetRoleByID(ctx, tenantID, dto.RoleId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("role not found")
		}
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}

	if dto.TeamId != 0 {
		_, err = u.teamRepo.GetTeam(ctx, tenantID, dto.TeamId, false)
		if err != nil {
			if ent.IsNotFound(err) {
				return v1.ErrorNotFound("team not found")
			}
			return v1.ErrorDatabaseQuery("get team failed: %v", err)
		}
	}

	err = u.repo.AssignRoles(ctx, tenantID, []data.AssignRoleDto{dto})
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	// TODO: remove check, with mock for QueueManager
	if u.qm != nil {
		u.qm.GetLocal(QueueRoleAssign).Pub(AssignRoleMessage{
			AssignRoleDto: dto,
			TenantID:      tenantID,
		})
	}

	return nil
}

func (u *AssignedRolesUsecase) UnassignRole(ctx context.Context, tenantID, assignID int64) error {
	assignedRole, err := u.repo.GetAssignedRoleById(ctx, tenantID, assignID)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("assigned role not found")
		}
		return v1.ErrorDatabaseQuery("get assigned role failed: %v", err)
	}

	err = u.repo.UnassignRole(ctx, assignedRole)
	if err != nil {
		return v1.ErrorDatabaseQuery("unassign role failed: %v", err)
	}

	// TODO: remove check, with mock for QueueManager
	if u.qm != nil {
		var resource *v1.Resource
		var teamID int64
		if assignedRole.ResourceID != nil {
			resource = &v1.Resource{
				Type: *assignedRole.ResourceType,
				Id:   *assignedRole.ResourceID,
			}

			if *assignedRole.ResourceType == data.RESOURCE_TYPE_TEAM { // for backward compatibility
				teamID = *assignedRole.ResourceID
			}
		}

		u.qm.GetLocal(QueueRoleUnassign).Pub(AssignRoleMessage{
			AssignRoleDto: data.AssignRoleDto{
				IdentityId: assignedRole.IdentityID,
				RoleId:     assignedRole.RoleID,
				TeamId:     teamID,
				Resource:   resource,
			},
			TenantID: tenantID,
		})
	}

	return nil
}

func (u *AssignedRolesUsecase) ListAssignedRoles(
	ctx context.Context,
	dto data.ListRolesDto,
) ([]*ent.ResourceAccess, error) {
	return u.repo.ListAssignedRoles(ctx, dto)
}
