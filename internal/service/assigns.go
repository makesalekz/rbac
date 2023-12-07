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

func (s *TeamIdentityRoleService) AssignRole(ctx context.Context, req *v1.AssignRoleRequest) (*utils_v1.EmptyReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can assign role to tenant identity

	err := s.uc.AssignRole(ctx, data.AssignRoleDto{
		TenantId:   claims.GetTenantId(),
		RoleId:     req.GetRoleId(),
		TeamId:     req.GetTeamId(),
		IdentityId: req.GetIdentityId(),
	})
	if err != nil {
		return nil, err
	}

	return &utils_v1.EmptyReply{}, nil
}

func (s *TeamIdentityRoleService) DeleteAssign(ctx context.Context, req *v1.AssignRequest) (*utils_v1.EmptyReply, error) {
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

func (s *TeamIdentityRoleService) ListAssigns(ctx context.Context, req *v1.ListAssignsRequest) (*v1.ListAssignsReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can get assgined roles

	identitiesIDs := []string{}
	teamsIDs := []int64{}

	if req.GetIdentityId() != "" {
		identitiesIDs = []string{req.GetIdentityId()}
	}

	if req.GetTeamId() != 0 {
		teamsIDs = []int64{req.GetTeamId()}
	}

	assignedRoles, err := s.uc.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    claims.GetTenantId(),
		IdentityIDs: identitiesIDs,
		TeamsIDs:    teamsIDs,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ListAssignsReply{
		Roles: assignedRolesReply(assignedRoles),
	}, nil
}

func assignedRoleReply(assignedRole *ent.TeamIdentityRole) *v1.AssignedRole {
	result := v1.AssignedRole{
		IdentityId: &assignedRole.IdentityID,
	}

	if assignedRole.Edges.Team != nil {
		result.Team = &v1.Team{
			Id:          assignedRole.Edges.Team.ID,
			Name:        assignedRole.Edges.Team.Name,
			Description: assignedRole.Edges.Team.Description,
		}
	}

	if assignedRole.Edges.Role != nil {
		result.Role = &v1.RoleReply{
			Id:          assignedRole.Edges.Role.ID,
			Name:        assignedRole.Edges.Role.Name,
			Description: assignedRole.Edges.Role.Description,
			IsSystem:    assignedRole.Edges.Role.IsSystem,
		}
	}

	return &result
}

func assignedRolesReply(assignedRoles []*ent.TeamIdentityRole) []*v1.AssignedRole {
	result := make([]*v1.AssignedRole, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		result[i] = assignedRoleReply(assignedRole)
	}
	return result
}
