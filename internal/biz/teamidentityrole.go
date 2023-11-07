package biz

import (
	"context"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	rbacv1 "rbac/api/rbac/v1"
	"rbac/ent"
	"rbac/internal/conf"
	"rbac/internal/data"
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
	// check jwt
	_, tenant, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, rbacv1.ErrorForbidden("forbidden")
	}
	// check role
	_, err := u.roleRepo.GetRoleById(ctx, dto.RoleId, tenant.TenantId)
	if err != nil {
		return nil, rbacv1.ErrorNotFound("role not found")
	}
	// check team
	_, err = u.teamRepo.GetTeam(ctx, dto.TeamId, false)
	if err != nil {
		return nil, rbacv1.ErrorNotFound("team not found")
	}
	teamIdentityRole, err := u.repo.AssignRole(ctx, dto)
	if err != nil {
		return nil, rbacv1.ErrorDatabaseQuery("assign role failed")
	}
	return teamIdentityRole, nil
}

func (u *TeamIdentityUsecase) DeleteIdentityRole(ctx context.Context, deleteDto data.DeleteRoleDto) error {
	_, tenant, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return rbacv1.ErrorUnauthorized("unauthorized")
	}
	_, err := u.repo.GetAssignedRoleById(ctx, deleteDto.AssignId, tenant.TenantId)
	if err != nil {
		return rbacv1.ErrorForbidden("forbidden")
	}
	return u.repo.DeleteIdentityRole(ctx, deleteDto)
}

func (u *TeamIdentityUsecase) ListIdentityRoles(ctx context.Context, dto data.ListIdentityRolesDto) ([]*ent.TeamIdentityRole, error) {
	_, tenant, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, rbacv1.ErrorUnauthorized("forbidden")
	}
	if dto.TenantId != tenant.TenantId {
		return nil, rbacv1.ErrorForbidden("forbidden")
	}
	return u.repo.ListIdentityRoles(ctx, dto)
}

func (u *TeamIdentityUsecase) ListTeamRoles(ctx context.Context, dto data.ListTeamRolesDto) ([]*ent.TeamIdentityRole, error) {
	_, tenant, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, rbacv1.ErrorUnauthorized("forbidden")
	}
	if dto.TenantId != tenant.TenantId {
		return nil, rbacv1.ErrorForbidden("forbidden")
	}
	return u.repo.ListTeamRoles(ctx, dto)
}

func (u *TeamIdentityUsecase) CheckPermissions(ctx context.Context, teamId int64, permissions []string) (map[string]*rbacv1.ListOfFields, error) {
	_, tenant, ok := u.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, rbacv1.ErrorUnauthorized("forbidden")
	}
	teamIdentityRoles, err := u.repo.ListTeamRoles(ctx, data.ListTeamRolesDto{
		TeamId:   teamId,
		TenantId: tenant.TenantId,
	})
	if err != nil {
		return nil, err
	}

	roleIds := make([]int64, 0)
	for _, teamIdentityRole := range teamIdentityRoles {
		roleIds = append(roleIds, teamIdentityRole.RoleID)
	}
	rolesPermissions, err := u.roleRepo.ListRolesPermissions(ctx, roleIds, tenant.TenantId, permissions)

	result := make(map[string]*rbacv1.ListOfFields)
	for _, rolePermission := range rolesPermissions {
		if _, ok := result[rolePermission.PermissionID]; !ok && !rolePermission.Deny {
			result[rolePermission.PermissionID] = &rbacv1.ListOfFields{
				Field: rolePermission.Fields,
			}
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
