package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTeamsUsecase_CreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	dto := data.TeamDto{
		TenantID: 1,
	}
	team := &ent.Team{
		ID: 1,
	}
	repo.EXPECT().CreateTeam(ctx, dto).Return(team, nil)

	team1, err := uc.CreateTeam(ctx, dto)
	require.NoError(t, err)
	require.Equal(t, team, team1)
}

func TestTeamsUsecase_CreateChildTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	parentID := int64(1)
	tenantID := int64(1)
	dto := data.TeamDto{
		ParentID:   parentID,
		TenantID:   tenantID,
		ParentsIDs: []int64{parentID},
	}
	dto2 := data.TeamDto{
		ParentID: 2,
		TenantID: tenantID,
	}
	dto3 := data.TeamDto{
		ParentID: parentID,
		TenantID: 2,
	}
	parentTeam := &ent.Team{
		ID:         parentID,
		TenantID:   tenantID,
		ParentsIds: nil,
	}
	team := &ent.Team{
		ID:       2,
		ParentID: &parentID,
		TenantID: tenantID,
	}

	repo.EXPECT().CreateTeam(ctx, dto).Return(team, nil)
	repo.EXPECT().GetTeam(ctx, parentID, tenantID, false).Return(parentTeam, nil)
	repo.EXPECT().GetTeam(ctx, gomock.Not(parentID), tenantID, false).Return(nil, &ent.NotFoundError{})
	repo.EXPECT().GetTeam(ctx, parentID, gomock.Not(tenantID), false).Return(nil, &ent.NotFoundError{})

	team1, err := uc.CreateTeam(ctx, dto)
	require.NoError(t, err)
	require.Equal(t, team, team1)

	team2, err := uc.CreateTeam(ctx, dto2)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("parent team not found"), err)
	require.Nil(t, team2)

	team3, err := uc.CreateTeam(ctx, dto3)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("parent team not found"), err)
	require.Nil(t, team3)
}

func TestTeamsUsecase_UpdateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	team := &ent.Team{
		ID: 1,
	}
	dto := data.TeamDto{
		TenantID: 1,
	}
	repo.EXPECT().UpdateTeam(ctx, team, dto).Return(team, nil)

	team1, err := uc.UpdateTeam(ctx, team, dto)
	require.NoError(t, err)
	require.Equal(t, team, team1)
}

func TestTeamsUsecase_DeleteTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	team := &ent.Team{
		ID: 1,
	}
	repo.EXPECT().DeleteTeam(ctx, team).Return(nil)

	err = uc.DeleteTeam(ctx, team)
	require.NoError(t, err)
}

func TestTeamsUsecase_GetTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	teamID := int64(1)
	team := &ent.Team{
		ID: 1,
	}
	repo.EXPECT().GetTeam(ctx, tenantID, teamID, false).Return(team, nil)
	repo.EXPECT().GetTeam(ctx, tenantID, gomock.Not(teamID), false).Return(nil, &ent.NotFoundError{})
	repo.EXPECT().GetTeam(ctx, gomock.Not(tenantID), teamID, false).Return(nil, &ent.NotFoundError{})

	team1, err := uc.GetTeam(ctx, tenantID, teamID, false)
	require.NoError(t, err)
	require.Equal(t, team, team1)

	team2, err := uc.GetTeam(ctx, tenantID, 2, false)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("team not found"), err)
	require.Nil(t, team2)

	team3, err := uc.GetTeam(ctx, 2, teamID, false)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("team not found"), err)
	require.Nil(t, team3)
}

func TestTeamsUsecase_ListTeams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewTeamsUsecase(repo)
	require.NoError(t, err)

	ctx := context.Background()
	filter := data.TeamsListFilter{}
	paginate := &utils_v1.PaginateRequest{}
	total := int32(2)
	teams := []*ent.Team{
		{
			ID:   1,
			Name: "team1",
		},
		{
			ID:   2,
			Name: "team2",
		},
	}
	teamsList := &biz.TeamsList{
		Teams: teams,
		Paginate: &utils_v1.PaginateReply{
			Total: &total,
		},
	}
	repo.EXPECT().ListTeams(ctx, filter, paginate).Return(teams, nil)
	repo.EXPECT().CountListTeams(ctx, filter).Return(int32(2), nil)

	teams1, err := uc.ListTeams(ctx, filter, nil)
	require.NoError(t, err)
	require.Equal(t, teamsList, teams1)
}
