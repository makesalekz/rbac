package data

import (
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
)

const ResourceTypeTeam = "team"

// ------- Assigns -----------

type AssignRoleDto struct {
	IdentityID string
	RoleID     int64
	TeamID     int64
	Resource   *v1.Resource
	Metadata   string
}

type ListRolesDto struct {
	TenantID       int64
	IdentityIDs    []string
	Resources      []*v1.Resource
	ResourceFilter []string
}

// ------- Permissions -------

type CreatePermissionDto struct {
	ID          string
	Name        string
	Description string
	GroupID     string
	AppID       string
	Fields      []string
}

func (dto CreatePermissionDto) Validate() error {
	if dto.ID == "" {
		return v1.ErrorBadRequest("empty id")
	}
	if dto.Name == "" {
		return v1.ErrorBadRequest("empty name")
	}
	if dto.GroupID == "" {
		return v1.ErrorBadRequest("empty group")
	}
	if dto.AppID == "" {
		return v1.ErrorBadRequest("empty app id")
	}
	// check if GroupId matches first part of Id
	// example: GroupId = "group1", Id = "group1.permission1"
	if dto.ID[:len(dto.GroupID)] != dto.GroupID {
		return v1.ErrorBadRequest("id must start with group id")
	}
	return nil
}

type UpdatePermissionDto struct {
	Name        string
	Description string
	Fields      []string
}

type FilterPermissions struct {
	AppsIDs    []string
	WithDenied bool
}

// ------- Roles -------------

type CreateRoleDto struct {
	Name        string
	Description string
	TenantID    int64
	IsSystem    bool
	Allow       []string
	Deny        []string
}

func (dto CreateRoleDto) Validate() error {
	if dto.Name == "" {
		return v1.ErrorBadRequest("empty name")
	}
	return nil
}

type UpdateRoleDto struct {
	Name        string
	Description string
	Allow       []string
	Deny        []string
}

type CreateRolePermissionDto struct {
	Deny   bool
	Fields []string
}

type FilterRolePermissions struct {
	AppIDs      []string
	TenantID    int64
	RoleIDs     []int64
	Permissions []string
	DeniedOnly  bool
}

// ------- Teams -------------

type TeamDto struct {
	TenantID    int64
	Name        string
	Description string
	ParentID    int64
	ParentsIDs  []int64
}

func (dto TeamDto) Validate() error {
	if dto.Name == "" {
		return v1.ErrorBadRequest("empty name")
	}
	return nil
}

type TeamsListFilter struct {
	TenantID int64
	ParentID int64
}
