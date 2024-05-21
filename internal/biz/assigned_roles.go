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
	TenantId int64
}

// AssignedRolesUsecase .
type AssignedRolesUsecase struct {
	log      *log.Helper
	repo     data.AssignedRolesRepo
	roleRepo data.RoleRepo
	teamRepo data.TeamsRepo
	queue    *u_nats.QueueManager
}

// NewAssignedRolesUsecase .
func NewAssignedRolesUsecase(
	logger log.Logger,
	repo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamRepo data.TeamsRepo,
	queue *u_nats.QueueManager,
) (*AssignedRolesUsecase, error) {
	return &AssignedRolesUsecase{
		log:      log.NewHelper(log.With(logger, "module", "biz/rbac")),
		repo:     repo,
		roleRepo: roleRepo,
		teamRepo: teamRepo,
		queue:    queue,
	}, nil
}

// AssignRoles assigns multiple roles to identities.
// This function checks if the role exists and if the team exists (if team id is not zero).
//
// Possible errors:
// - rbac.ErrorNotFound: role or team not found
// - rbac.ErrorDatabaseQuery: failed to get role or team
// - rbac.ErrorAlreadyExists: role already assigned
// - rbac.ErrorBadRequest: there is no such teamId
func (u *AssignedRolesUsecase) AssignRoles(ctx context.Context, tenantId int64, dtos []data.AssignRoleDto) error {
	roleIds := data.ExtractSlice(dtos, func(e data.AssignRoleDto) (int64, bool) { return e.RoleId, true })
	// Get roles by ids
	_, err := u.roleRepo.GetRolesById(ctx, tenantId, roleIds)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("role not found")
		}
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}

	teamsIds := data.ExtractSlice(dtos, func(e data.AssignRoleDto) (int64, bool) {
		if e.TeamId != 0 {
			return e.TeamId, true
		}
		return 0, false
	})

	// If there are team ids, get teams by ids, then check if the returned ids are equal to the input ids
	if len(teamsIds) > 0 {
		teams, err := u.teamRepo.GetTeams(ctx, tenantId, teamsIds)
		if err != nil {
			if len(teams) == 0 {
				return v1.ErrorNotFound("teams not found")
			}
			return v1.ErrorDatabaseQuery("get team failed: %v", err)
		}

		returnedIds := map[int64]struct{}{}
		for _, team := range teams {
			returnedIds[team.ID] = struct{}{}
		}

		for _, teamId := range teamsIds {
			_, ok := returnedIds[teamId]
			if !ok {
				return v1.ErrorBadRequest("there is no such teamId: %d", teamId)
			}
		}
	}

	// Assign roles
	err = u.repo.AssignRoles(ctx, tenantId, dtos)
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	for _, dto := range dtos {
		if u.queue == nil {
			continue
		}

		u.queue.GetLocal(QueueRoleAssign).Pub(AssignRoleMessage{
			AssignRoleDto: dto,
			TenantId:      tenantId,
		})
	}

	return nil
}

func (u *AssignedRolesUsecase) AssignRole(ctx context.Context, tenantId int64, dto data.AssignRoleDto) error {
	_, err := u.roleRepo.GetRoleById(ctx, tenantId, dto.RoleId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("role not found")
		}
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}

	if dto.TeamId != 0 {
		_, err = u.teamRepo.GetTeam(ctx, tenantId, dto.TeamId, false)
		if err != nil {
			if ent.IsNotFound(err) {
				return v1.ErrorNotFound("team not found")
			}
			return v1.ErrorDatabaseQuery("get team failed: %v", err)
		}
	}

	err = u.repo.AssignRoles(ctx, tenantId, []data.AssignRoleDto{dto})
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	if u.queue != nil {
		u.queue.GetLocal(QueueRoleAssign).Pub(AssignRoleMessage{
			AssignRoleDto: dto,
			TenantId:      tenantId,
		})
	}

	return nil
}

func (u *AssignedRolesUsecase) UnassignRole(ctx context.Context, tenantId, assignId int64) error {
	assignedRole, err := u.repo.GetAssignedRoleById(ctx, tenantId, assignId)
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

	if u.queue != nil {
		u.queue.GetLocal(QueueRoleUnassign).Pub(AssignRoleMessage{
			AssignRoleDto: data.AssignRoleDto{
				IdentityId: assignedRole.IdentityID,
				TeamId:     *assignedRole.TeamID,
				RoleId:     assignedRole.RoleID,
			},
			TenantId: tenantId,
		})
	}

	return nil
}

func (u *AssignedRolesUsecase) ListIdentityRoles(ctx context.Context, tenantId int64, identityId string) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: []string{identityId},
	})
}

func (u *AssignedRolesUsecase) ListAssignedRoles(ctx context.Context, dto data.ListRolesDto) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListAssignedRoles(ctx, dto)
}
