package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

type TeamIdentityRoleService struct {
	v1.UnimplementedTeamIdentityRoleServer

	jwt *jwt.JwtProcessor
	ru  *biz.RolesUsecase
	tu  *biz.TeamsUsecase
	uc  *biz.TeamIdentityUsecase
}

func NewTeamIdentityRoleService(
	jwt *jwt.JwtProcessor,
	ru *biz.RolesUsecase,
	tu *biz.TeamsUsecase,
	uc *biz.TeamIdentityUsecase,
) *TeamIdentityRoleService {
	return &TeamIdentityRoleService{
		jwt: jwt,
		ru:  ru,
		tu:  tu,
		uc:  uc,
	}
}

func (s *TeamIdentityRoleService) AssignRole(ctx context.Context, req *v1.AssignRoleRequest) (*v1.TeamIdentityRoleReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can assign role to tenant identity

	_, err := s.ru.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	if req.TeamId != 0 {
		_, err = s.tu.GetTeam(ctx, claims.GetTenantId(), req.TeamId, false)
		if err != nil {
			return nil, err
		}
	}

	assignedRole, err := s.uc.AssignRole(ctx, data.AssignRoleDto{
		RoleId:     req.RoleId,
		TeamId:     req.TeamId,
		IdentityId: req.IdentityId,
		TenantId:   claims.GetTenantId(),
	})
	if err != nil {
		return nil, err
	}
	return assignedRoleReply(assignedRole), nil
}

func (s *TeamIdentityRoleService) DeleteRole(ctx context.Context, req *v1.DeleteRequest) (*utils_v1.EmptyReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can delete role

	err := s.uc.DeleteIdentityRole(ctx, claims.GetTenantId(), req.GetAssignId())
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *TeamIdentityRoleService) ListIdentityRoles(ctx context.Context, req *v1.ListIdentityRolesRequest) (*v1.ListIdentityRolesReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can get identity roles

	assignedRoles, err := s.uc.ListIdentityRoles(ctx, claims.GetTenantId(), req.GetIdentityId())
	if err != nil {
		return nil, err
	}

	return &v1.ListIdentityRolesReply{
		Roles: assignedRolesReply(assignedRoles),
	}, nil
}

func (s *TeamIdentityRoleService) ListTeamRoles(ctx context.Context, req *v1.ListTeamRolesRequest) (*v1.ListIdentityRolesReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can get assgined roles

	assignedRoles, err := s.uc.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId: claims.GetTenantId(),
		TeamId:   req.GetTeamId(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.ListIdentityRolesReply{
		Roles: assignedRolesReply(assignedRoles),
	}, nil
}

func assignedRoleReply(assignedRole *ent.TeamIdentityRole) *v1.TeamIdentityRoleReply {
	return &v1.TeamIdentityRoleReply{
		IdentityId: &assignedRole.IdentityID,
		Team: &v1.Team{
			Id:          assignedRole.Edges.Team.ID,
			Name:        assignedRole.Edges.Team.Name,
			Description: assignedRole.Edges.Team.Description,
		},
		Role: &v1.RoleReply{
			Id:          assignedRole.RoleID,
			Name:        assignedRole.Edges.Role.Name,
			Description: assignedRole.Edges.Role.Description,
			IsSystem:    assignedRole.Edges.Role.IsSystem,
		},
	}
}

func assignedRolesReply(assignedRoles []*ent.TeamIdentityRole) []*v1.TeamIdentityRoleReply {
	result := make([]*v1.TeamIdentityRoleReply, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		result[i] = assignedRoleReply(assignedRole)
	}
	return result
}
