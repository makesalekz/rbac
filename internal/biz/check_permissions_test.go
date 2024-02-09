package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCheckPermissionsUsecase_CheckPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identities := []string{"identity1", "identity2"}
	teamId := int64(1)
	permissions := []string{"permission.one", "permission.two"}

	team := &ent.Team{
		ID:         teamId,
		TenantID:   tenantId,
		Name:       "testName",
		ParentsIds: &pgtype.Int8Array{},
	}
	teamRepo.EXPECT().GetTeam(ctx, tenantId, teamId, false).Return(team, nil)

	assignedRoles := []*ent.TeamIdentityRole{
		{
			ID:         1,
			TenantID:   tenantId,
			IdentityID: "identity1",
			RoleID:     1,
		},
		{
			ID:         2,
			TenantID:   tenantId,
			IdentityID: "identity2",
			RoleID:     2,
		},
	}
	listRolesDto := data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: identities,
		TeamsIDs:    []int64{teamId},
	}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, listRolesDto).Return(assignedRoles, nil)

	filterRolePermissions := data.FilterRolePermissions{
		TenantId:    tenantId,
		RolesIds:    []int64{1, 2},
		Permissions: permissions,
	}
	rolesPermissions := []*ent.RolePermission{
		{
			ID:           1,
			TenantID:     tenantId,
			RoleID:       1,
			PermissionID: "permission.one",
			Fields:       []string{"field1", "field2"},
		},
		{
			ID:           2,
			TenantID:     tenantId,
			RoleID:       2,
			PermissionID: "permission.two",
			Fields:       []string{"field3", "field4"},
		},
	}
	roleRepo.EXPECT().ListRolesPermissions(ctx, filterRolePermissions).Return(rolesPermissions, nil)

	permissionsMap, err := uc.CheckPermissions(ctx, tenantId, identities, teamId, permissions)
	require.NoError(t, err)
	require.Len(t, permissionsMap, 2)

	require.Equal(t, &v1.ListOfFields{
		Fields: []string{"field1", "field2"},
	}, permissionsMap["permission.one"])

	require.Equal(t, &v1.ListOfFields{
		Fields: []string{"field3", "field4"},
	}, permissionsMap["permission.two"])
}
