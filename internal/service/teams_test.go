package service_test

import (
	"testing"
	"time"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/rbac/internal/service"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

func createTeamsService(
	t *testing.T,
	assignedRepo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamsRepo data.TeamsRepo,
) *service.TeamsService {
	tu, err := biz.NewTeamsUsecase(teamsRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	return service.NewTeamsService(sh, tu)
}

// ------------------ CreateTeam ------------------------

func TestRolesService_CreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.CreateTeamRequest{
		Name:        "testName",
		Description: "testDescription",
	}
	dto := data.TeamDto{
		TenantID:    tenantID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
	team := &ent.Team{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teamsRepo.EXPECT().CreateTeam(ctx, dto).Return(team, nil)

	expect := &v1.Team{
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
	}

	reply, err := service.CreateTeam(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetTeam())
}

func TestRolesService_CreateTeamWithParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	parentIDs := []pgtype.Int8{
		{Int: 1, Status: pgtype.Present},
		{Int: 2, Status: pgtype.Present},
	}
	parentTeam := &ent.Team{
		ID:          3,
		Name:        "parentName",
		Description: "parentDescription",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ParentsIds:  &pgtype.Int8Array{Elements: parentIDs, Status: pgtype.Present},
	}
	req := &v1.CreateTeamRequest{
		Name:        "testName",
		Description: "testDescription",
		ParentId:    parentTeam.ID,
	}
	dto := data.TeamDto{
		TenantID:    tenantID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ParentID:    req.GetParentId(),
		ParentsIDs:  []int64{1, 2, 3},
	}
	teamParentIDs := &pgtype.Int8Array{
		Elements: append(parentIDs, pgtype.Int8{Int: 3, Status: pgtype.Present}),
		Status:   pgtype.Present,
	}
	team := &ent.Team{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ParentID:    &parentTeam.ID,
		ParentsIds:  teamParentIDs,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, parentTeam.ID, false).Return(parentTeam, nil)
	teamsRepo.EXPECT().CreateTeam(ctx, dto).Return(team, nil)

	expect := &v1.Team{
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
		ParentId:    &parentTeam.ID,
		ParentsIds:  []int64{1, 2, 3},
	}

	reply, err := service.CreateTeam(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetTeam())
}

func TestRolesService_CreateTeamEmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.CreateTeamRequest{
		Description: "testDescription",
	}

	reply, err := service.CreateTeam(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty name"), err)
	require.Nil(t, reply)
}

// ------------------ UpdateTeam ------------------------

func TestRolesService_UpdateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.UpdateTeamRequest{
		TeamId:      1,
		Name:        "testName",
		Description: "testDescription",
	}
	dto := data.TeamDto{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
	team := &ent.Team{
		ID:          req.GetTeamId(),
		Name:        "oldName",
		Description: "oldDescription",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	updatedTeam := &ent.Team{
		ID:          req.GetTeamId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(team, nil)
	teamsRepo.EXPECT().UpdateTeam(ctx, team, dto).Return(updatedTeam, nil)

	expect := &v1.Team{
		Id:          updatedTeam.ID,
		Name:        updatedTeam.Name,
		Description: updatedTeam.Description,
		CreatedAt:   updatedTeam.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTeam.UpdatedAt.Format(time.RFC3339),
	}

	reply, err := service.UpdateTeam(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetTeam())
}

func TestRolesService_UpdateTeamEmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.UpdateTeamRequest{
		Name:        "testName",
		Description: "testDescription",
	}

	reply, err := service.UpdateTeam(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty team id"), err)
	require.Nil(t, reply)
}

// ------------------ DeleteTeam ------------------------

func TestRolesService_DeleteTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.TeamRequest{
		TeamId: 1,
	}
	team := &ent.Team{
		ID:          req.GetTeamId(),
		Name:        "oldName",
		Description: "oldDescription",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(team, nil)
	teamsRepo.EXPECT().DeleteTeam(ctx, team).Return(nil)

	_, err := service.DeleteTeam(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_DeleteTeamEmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.TeamRequest{}

	reply, err := service.DeleteTeam(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty team id"), err)
	require.Nil(t, reply)
}

// ------------------ GetTeam ---------------------------

func TestRolesService_GetTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.TeamRequest{
		TeamId: 1,
	}
	team := &ent.Team{
		ID:          req.GetTeamId(),
		Name:        "oldName",
		Description: "oldDescription",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teamsRepo.EXPECT().GetTeam(ctx, tenantID, req.GetTeamId(), false).Return(team, nil)

	expect := &v1.Team{
		Id:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
	}

	reply, err := service.GetTeam(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetTeam())
}

func TestRolesService_GetTeamEmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.TeamRequest{}

	reply, err := service.GetTeam(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty team id"), err)
	require.Nil(t, reply)
}

// ------------------ ListTeams -------------------------

func TestRolesService_ListTeams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createTeamsService(t, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.ListTeamsRequest{
		ParentId: 123,
		Paginate: &utils_v1.PaginateRequest{
			Page:  1,
			Limit: 10,
		},
	}
	filter := data.TeamsListFilter{
		TenantID: tenantID,
		ParentID: req.GetParentId(),
	}
	team := &ent.Team{
		ID:          1,
		Name:        "testName",
		Description: "testDescription",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	teams := []*ent.Team{
		team,
		{
			ID:          2,
			Name:        "testName2",
			Description: "testDescription2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	teamsRepo.EXPECT().ListTeams(ctx, filter, req.GetPaginate()).Return(teams, nil)
	teamsRepo.EXPECT().CountListTeams(ctx, filter).Return(int32(2), nil)

	expect := &v1.Team{
		Id:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
	}

	reply, err := service.ListTeams(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetTeams(), 2)
	require.Equal(t, expect, reply.GetTeams()[0])
}
