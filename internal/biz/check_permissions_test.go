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
	tenantID := int64(1)
	identities := []string{"identity1", "identity2"}
	teamID := int64(1)
	permissions := []string{"permission.one", "permission.two"}

	team := &ent.Team{
		ID:         teamID,
		TenantID:   tenantID,
		Name:       "testName",
		ParentsIds: nil,
	}
	teamRepo.EXPECT().GetTeams(ctx, tenantID, []int64{teamID}).Return([]*ent.Team{team}, nil)

	assignedRoles := []*ent.ResourceAccess{
		{
			ID:         1,
			TenantID:   tenantID,
			IdentityID: "identity1",
			RoleID:     1,
		},
		{
			ID:         2,
			TenantID:   tenantID,
			IdentityID: "identity2",
			RoleID:     2,
		},
	}
	listRolesDto := data.ListRolesDto{
		TenantID:    tenantID,
		IdentityIDs: identities,
		Resources: []*v1.Resource{
			{Id: teamID, Type: data.ResourceTypeTeam},
		},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, listRolesDto).Return(assignedRoles, nil)

	filterRolePermissions := data.FilterRolePermissions{
		TenantID:    tenantID,
		RoleIDs:     []int64{1, 2},
		Permissions: permissions,
	}
	rolesPermissions := []*ent.RolePermission{
		{
			ID:           1,
			TenantID:     tenantID,
			RoleID:       1,
			PermissionID: "permission.one",
			Fields:       []string{"field1", "field2"},
		},
		{
			ID:           2,
			TenantID:     tenantID,
			RoleID:       2,
			PermissionID: "permission.two",
			Fields:       []string{"field3", "field4"},
		},
	}
	roleRepo.EXPECT().ListRolesPermissions(ctx, filterRolePermissions).Return(rolesPermissions, nil)

	resources := []*v1.Resource{{Id: teamID, Type: data.ResourceTypeTeam}}

	permissionsMap, err := uc.CheckPermissions(ctx, tenantID, identities, permissions, resources)
	require.NoError(t, err)
	require.Len(t, permissionsMap, 2)

	require.Equal(t, &v1.ListOfFields{
		Fields: []string{"field1", "field2"},
	}, permissionsMap["permission.one"])

	require.Equal(t, &v1.ListOfFields{
		Fields: []string{"field3", "field4"},
	}, permissionsMap["permission.two"])
}
