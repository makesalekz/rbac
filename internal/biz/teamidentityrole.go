package biz

import (
	"context"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/conf"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

// TeamIdentityUsecase .
type TeamIdentityUsecase struct {
	conf      *conf.Bootstrap
	log       *log.Helper
	discovery *consul.Registry
	jwt       *data.JwtProcessor
	repo      data.TeamIdentityRoleRepo
	roleRepo  data.RoleRepo
	teamRepo  data.TeamsRepo
}

func (u *TeamIdentityUsecase) AssignRole(ctx context.Context, dto data.AssignRoleDto) (*ent.TeamIdentityRole, error) {
	claims, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorForbidden("forbidden")
	}
	if dto.TenantId != claims.TenantId {
		return nil, v1.ErrorForbidden("forbidden")
	}
	// todo checkPermissions can assign role to tenant identity
	_, err := u.roleRepo.GetRoleById(ctx, dto.RoleId, claims.TenantId)
	if err != nil {
		return nil, v1.ErrorNotFound("role not found")
	}
	_, err = u.teamRepo.GetTeam(ctx, dto.TeamId, claims.TenantId, false)
	if err != nil {
		return nil, v1.ErrorNotFound("team not found")
	}

	teamIdentityRole, err := u.repo.AssignRole(ctx, dto)
	if err != nil {
		return nil, v1.ErrorDatabaseQuery("assign role failed")
	}
	return teamIdentityRole, nil
}

func (u *TeamIdentityUsecase) DeleteIdentityRole(ctx context.Context, deleteDto data.DeleteRoleDto) error {
	claims, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("unauthorized")
	}
	// todo checkPermissions can delete identity role
	_, err := u.repo.GetAssignedRoleById(ctx, deleteDto.AssignId, claims.TenantId)
	if err != nil {
		return v1.ErrorForbidden("forbidden")
	}
	return u.repo.DeleteIdentityRole(ctx, deleteDto)
}

func (u *TeamIdentityUsecase) ListIdentityRoles(ctx context.Context, dto data.ListIdentityRolesDto) ([]*ent.TeamIdentityRole, error) {
	claims, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("forbidden")
	}
	if dto.TenantId != claims.TenantId {
		return nil, v1.ErrorForbidden("forbidden")
	}
	return u.repo.ListIdentityRoles(ctx, dto)
}

func (u *TeamIdentityUsecase) ListTeamRoles(ctx context.Context, dto data.ListTeamRolesDto) ([]*ent.TeamIdentityRole, error) {
	claims, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("forbidden")
	}
	if dto.TenantId != claims.TenantId {
		return nil, v1.ErrorForbidden("forbidden")
	}
	return u.repo.ListTeamRoles(ctx, dto)
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

func (u *TeamIdentityUsecase) CheckPermissions(ctx context.Context, teamId int64, permissions []string) (map[string]*v1.ListOfFields, error) {
	claims, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("forbidden")
	}
	teamIdentityRoles, err := u.repo.ListTeamRoles(ctx, data.ListTeamRolesDto{
		TeamId:      teamId,
		TenantId:    claims.TenantId,
		IdentityIDs: claims.GetIdentities(),
	})
	if err != nil {
		return nil, err
	}

	roleIds := make([]int64, 0)
	for _, teamIdentityRole := range teamIdentityRoles {
		roleIds = append(roleIds, teamIdentityRole.RoleID)
	}
	rolesPermissions, err := u.roleRepo.ListRolesPermissions(ctx, roleIds, claims.TenantId, permissions)
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

// NewTeamIdentityUsecase .
func NewTeamIdentityUsecase(logger log.Logger, c *data.Config, jwt *data.JwtProcessor, repo data.TeamIdentityRoleRepo, roleRepo data.RoleRepo, teamRepo data.TeamsRepo) (*TeamIdentityUsecase, error) {
	return &TeamIdentityUsecase{
		conf:      c.Bootstrap,
		log:       log.NewHelper(logger),
		discovery: c.GetRegistry(),
		jwt:       jwt,
		repo:      repo,
		roleRepo:  roleRepo,
		teamRepo:  teamRepo,
	}, nil
}
