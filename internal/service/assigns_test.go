package service_test

import (
	"errors"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/rbac/internal/service"
	u_nats "gitlab.calendaria.team/services/utils/v1/nats"
	u_zap "gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func createAssignsService(
	t *testing.T,
	assignedRepo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamsRepo data.TeamsRepo,
) *service.AssignsService {
	logger := u_zap.NewZapLogger(true)

	var qm *u_nats.QueueManager

	ru, err := biz.NewRolesUsecase(logger, roleRepo)
	require.NoError(t, err)

	au, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamsRepo, qm)
	require.NoError(t, err)

	tu, err := biz.NewTeamsUsecase(teamsRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	return service.NewAssignsService(ru, tu, au, sh)
}

// ------------------ AssignRoles ------------------------

func TestRolesService_AssignRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRolesRequest{
		Assigns: []*v1.AssignRoleRequest{
			{
				IdentityId: "1234",
				RoleId:     11,
				TeamId:     22,
			},
			{
				RoleId: 12,
			},
		},
	}
	roles := []*ent.Role{
		{ID: 11},
		{ID: 12},
	}
	teams := []*ent.Team{
		{ID: 22},
	}
	dtos := []data.AssignRoleDto{
		{IdentityID: "1234", RoleID: 11, TeamID: 22, Resource: &v1.Resource{Type: data.RESOURCE_TYPE_TEAM, Id: 22}},
		{RoleID: 12},
	}
	roleRepo.EXPECT().GetRolesByID(ctx, tenantID, []int64{11, 12}).Return(roles, nil)
	teamsRepo.EXPECT().GetTeams(ctx, tenantID, []int64{22}).Return(teams, nil)
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, dtos).Return(nil)

	_, err := service.AssignRoles(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_AssignRolesEmptyRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRolesRequest{
		Assigns: []*v1.AssignRoleRequest{
			{
				IdentityId: "1234",
				TeamId:     22,
			},
			{
				RoleId: 12,
			},
		},
	}

	reply, err := service.AssignRoles(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRolesInvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRolesRequest{
		Assigns: []*v1.AssignRoleRequest{
			{
				IdentityId: "1234",
				RoleId:     11,
				TeamId:     22,
			},
			{
				RoleId: 12,
			},
			{
				RoleId: 13,
			},
		},
	}
	roles := []*ent.Role{
		{ID: 11},
	}
	roleRepo.EXPECT().GetRolesByID(ctx, tenantID, []int64{11, 12, 13}).Return(roles, nil)

	reply, err := service.AssignRoles(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("invalid role ids [12 13]"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRolesInvalidTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRolesRequest{
		Assigns: []*v1.AssignRoleRequest{
			{
				IdentityId: "1234",
				RoleId:     11,
				TeamId:     22,
			},
			{
				RoleId: 12,
				TeamId: 23,
			},
			{
				RoleId: 13,
				TeamId: 24,
			},
		},
	}
	roles := []*ent.Role{
		{ID: 11},
		{ID: 12},
		{ID: 13},
	}
	teams := []*ent.Team{
		{ID: 22},
	}
	roleRepo.EXPECT().GetRolesByID(ctx, tenantID, []int64{11, 12, 13}).Return(roles, nil)
	teamsRepo.EXPECT().GetTeams(ctx, tenantID, []int64{22, 23, 24}).Return(teams, nil)

	reply, err := service.AssignRoles(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("invalid team ids [23 24]"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRolesAlreadyAssigned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRolesRequest{
		Assigns: []*v1.AssignRoleRequest{
			{
				IdentityId: "1234",
				RoleId:     11,
				TeamId:     22,
			},
			{
				RoleId: 12,
			},
		},
	}
	roles := []*ent.Role{
		{ID: 11},
		{ID: 12},
	}
	teams := []*ent.Team{
		{ID: 22},
	}
	dtos := []data.AssignRoleDto{
		{IdentityID: "1234", RoleID: 11, TeamID: 22, Resource: &v1.Resource{Type: data.RESOURCE_TYPE_TEAM, Id: 22}},
		{RoleID: 12},
	}
	roleRepo.EXPECT().GetRolesByID(ctx, tenantID, []int64{11, 12}).Return(roles, nil)
	teamsRepo.EXPECT().GetTeams(ctx, tenantID, []int64{22}).Return(teams, nil)
	e := ent.NewConstraintError("id exists", errors.New("id exists"))
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, dtos).Return(e)

	reply, err := service.AssignRoles(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorAlreadyExists("role already assigned"), err)
	require.Nil(t, reply)
}

// ------------------ AssignRole -------------------------

func TestRolesService_AssignRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRoleRequest{
		IdentityId: "1234",
		RoleId:     11,
		TeamId:     22,
	}
	role := &ent.Role{ID: req.GetRoleId()}
	team := &ent.Team{ID: req.GetTeamId()}
	dtos := []data.AssignRoleDto{
		{
			IdentityID: req.GetIdentityId(),
			RoleID:     req.GetRoleId(),
			TeamID:     req.GetTeamId(),
			Resource:   &v1.Resource{Type: data.RESOURCE_TYPE_TEAM, Id: req.GetTeamId()},
		},
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, req.GetRoleId()).Return(role, nil)
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(team, nil)
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, dtos).Return(nil)

	_, err := service.AssignRole(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_AssignRoleEmptyRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRoleRequest{
		IdentityId: "1234",
		TeamId:     22,
	}

	reply, err := service.AssignRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRoleNotFoundRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRoleRequest{
		IdentityId: "1234",
		RoleId:     11,
		TeamId:     22,
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, req.GetRoleId()).Return(nil, ent.NewNotFoundError("not found"))

	reply, err := service.AssignRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRoleNotFoundTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRoleRequest{
		IdentityId: "1234",
		RoleId:     11,
		TeamId:     22,
	}
	role := &ent.Role{ID: req.GetRoleId()}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, req.GetRoleId()).Return(role, nil)
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(nil, ent.NewNotFoundError("not found"))

	reply, err := service.AssignRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("team not found"), err)
	require.Nil(t, reply)
}

func TestRolesService_AssignRoleAlreadyAssigned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRoleRequest{
		IdentityId: "1234",
		RoleId:     11,
		TeamId:     22,
	}
	role := &ent.Role{ID: req.GetRoleId()}
	team := &ent.Team{ID: req.GetTeamId()}
	dtos := []data.AssignRoleDto{
		{
			IdentityID: req.GetIdentityId(),
			RoleID:     req.GetRoleId(),
			TeamID:     req.GetTeamId(),
			Resource:   &v1.Resource{Type: data.RESOURCE_TYPE_TEAM, Id: req.GetTeamId()},
		},
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, req.GetRoleId()).Return(role, nil)
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(team, nil)
	e := ent.NewConstraintError("id exists", errors.New("id exists"))
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, dtos).Return(e)

	reply, err := service.AssignRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorAlreadyExists("role already assigned"), err)
	require.Nil(t, reply)
}

// ------------------ UnassignRole -----------------------

func TestRolesService_UnassignRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRequest{
		AssignId: 1234,
	}
	assignedRole := &ent.ResourceAccess{
		ID: req.GetAssignId(),
	}
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantID, req.GetAssignId()).Return(assignedRole, nil)
	assignedRepo.EXPECT().UnassignRole(ctx, assignedRole).Return(nil)

	_, err := service.UnassignRole(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_UnassignRoleNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AssignRequest{
		AssignId: 1234,
	}
	e := ent.NewNotFoundError("not found")
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantID, req.GetAssignId()).Return(nil, e)

	reply, err := service.UnassignRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("assigned role not found"), err)
	require.Nil(t, reply)
}

// ------------------ ListAssigns ------------------------

func TestRolesService_ListAssigns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createAssignsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	resource := &v1.Resource{Type: "test", Id: 222}
	role0 := &ent.Role{ID: 11, Name: "role0"}
	role1 := &ent.Role{ID: 12, Name: "role1"}
	assignedRole0 := &ent.ResourceAccess{
		ID:           123,
		IdentityID:   "1234",
		RoleID:       role0.ID,
		ResourceType: &resource.Type,
		ResourceID:   &resource.Id,
		Edges:        ent.ResourceAccessEdges{Role: role0},
	}
	assignedRole1 := &ent.ResourceAccess{
		ID:         124,
		IdentityID: assignedRole0.IdentityID,
		RoleID:     role1.ID,
		Edges:      ent.ResourceAccessEdges{Role: role1},
	}

	req := &v1.ListAssignsRequest{
		IdentityIds: []string{assignedRole0.IdentityID},
		Resources:   []*v1.Resource{resource},
	}
	dto := data.ListRolesDto{
		TenantID:       tenantID,
		IdentityIDs:    req.GetIdentityIds(),
		Resources:      req.GetResources(),
		ResourceFilter: req.GetResourceTypes(),
	}
	assignedRoles := []*ent.ResourceAccess{
		assignedRole0,
		assignedRole1,
	}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, dto).Return(assignedRoles, nil)

	expect0 := &v1.AssignedRole{
		AssignId:   assignedRole0.ID,
		IdentityId: &assignedRole0.IdentityID,
		Role: &v1.Role{
			Id:   role0.ID,
			Name: role0.Name,
		},
		Resource: resource,
	}
	expect1 := &v1.AssignedRole{
		AssignId:   assignedRole1.ID,
		IdentityId: &assignedRole1.IdentityID,
		Role: &v1.Role{
			Id:   role1.ID,
			Name: role1.Name,
		},
	}

	reply, err := service.ListAssigns(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetRoles(), 2)
	require.Equal(t, expect0, reply.GetRoles()[0])
	require.Equal(t, expect1, reply.GetRoles()[1])
}
