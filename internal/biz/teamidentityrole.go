package biz

import (
	"context"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/conf"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/utils/v1/config"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

// TeamIdentityUsecase .
type TeamIdentityUsecase struct {
	conf      *conf.Bootstrap
	log       *log.Helper
	discovery *consul.Registry
	jwt       *jwt.JwtProcessor
	repo      data.TeamIdentityRoleRepo
	roleRepo  data.RoleRepo
	teamRepo  data.TeamsRepo
}

func (u *TeamIdentityUsecase) AssignRole(ctx context.Context, dto data.AssignRoleDto) error {
	_, err := u.roleRepo.GetRoleById(ctx, dto.TenantId, dto.RoleId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("role not found")
		}
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}

	if dto.TeamId != 0 {
		_, err = u.teamRepo.GetTeam(ctx, dto.TenantId, dto.TeamId, false)
		if err != nil {
			if ent.IsNotFound(err) {
				return v1.ErrorNotFound("team not found")
			}
			return v1.ErrorDatabaseQuery("get team failed: %v", err)
		}
	}

	err = u.repo.AssignRole(ctx, dto)
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	return nil
}

func (u *TeamIdentityUsecase) DeleteIdentityRole(ctx context.Context, tenantId, assignId int64) error {
	assignedRole, err := u.repo.GetAssignedRoleById(ctx, tenantId, assignId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("assgined role not found")
		}
		return v1.ErrorDatabaseQuery("get assgined role failed: %v", err)
	}

	return u.repo.DeleteIdentityRole(ctx, assignedRole)
}

func (u *TeamIdentityUsecase) ListIdentityRoles(ctx context.Context, tenantId int64, identityId string) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListRoles(ctx, data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: []string{identityId},
	})
}

func (u *TeamIdentityUsecase) ListAssignedRoles(ctx context.Context, dto data.ListRolesDto) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListRoles(ctx, dto)
}

func (u *TeamIdentityUsecase) CheckPermissions(ctx context.Context, teamId int64, permissions []string) (map[string]*v1.ListOfFields, error) {
	claims, ok := u.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}

	u.log.Debugf("check permissions: %+v", claims)

	var teamsIds []int64
	if teamId != 0 {
		team, err := u.teamRepo.GetTeam(ctx, claims.GetTenantId(), teamId, false)
		if err != nil {
			return nil, err
		}
		team.ParentsIds.AssignTo(&teamsIds)
		teamsIds = append(teamsIds, team.ID)
	}

	assignedRoles, err := u.repo.ListRoles(ctx, data.ListRolesDto{
		TenantId:    claims.GetTenantId(),
		IdentityIDs: claims.GetIdentities(),
		TeamsIDs:    teamsIds,
	})
	if err != nil {
		return nil, err
	}

	rolesIds := make([]int64, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		rolesIds[i] = assignedRole.RoleID
	}

	rolesPermissions, err := u.roleRepo.ListRolesPermissions(ctx, rolesIds, claims.GetTenantId(), permissions)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*v1.ListOfFields)
	for _, rolePermission := range rolesPermissions {
		if _, ok := result[rolePermission.PermissionID]; !ok {
			result[rolePermission.PermissionID] = &v1.ListOfFields{
				Fields: rolePermission.Fields,
			}
		} else {
			result[rolePermission.PermissionID].Fields = mergeFields(result[rolePermission.PermissionID].Fields, rolePermission.Fields)
		}
	}

	for _, rolePermission := range rolesPermissions {
		if rolePermission.Deny {
			delete(result, rolePermission.PermissionID)
		}
	}

	return result, nil
}

func (u *TeamIdentityUsecase) HasPermission(ctx context.Context, permission string) (*v1.ListOfFields, error) {
	permissionsMap, err := u.CheckPermissions(ctx, 0, []string{permission})
	if err != nil {
		return nil, err
	}

	if len(permissionsMap) == 0 {
		return nil, nil
	}

	return permissionsMap[permission], nil
}

// NewTeamIdentityUsecase .
func NewTeamIdentityUsecase(
	conf *conf.Bootstrap,
	logger log.Logger,
	c *config.Config,
	jwt *jwt.JwtProcessor,
	repo data.TeamIdentityRoleRepo,
	roleRepo data.RoleRepo,
	teamRepo data.TeamsRepo,
) (*TeamIdentityUsecase, error) {
	return &TeamIdentityUsecase{
		conf:      conf,
		log:       log.NewHelper(logger),
		discovery: c.GetRegistry(),
		jwt:       jwt,
		repo:      repo,
		roleRepo:  roleRepo,
		teamRepo:  teamRepo,
	}, nil
}

func mergeFields(fields1 []string, fields2 []string) []string {
	if len(fields1) > len(fields2) {
		fields1, fields2 = fields2, fields1
	}
	for _, field := range fields2 {
		if !contains(fields1, field) {
			fields1 = append(fields1, field)
		}
	}
	return fields1
}

func contains(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}
