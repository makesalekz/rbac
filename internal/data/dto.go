package data

import (
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
)

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
	TenantID    int64
	RolesIDs    []int64
	Permissions []string
	DeniedOnly  bool
}
