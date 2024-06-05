package data

import (
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
)

type CreatePermissionDto struct {
	Id          string
	Name        string
	Description string
	GroupId     string
	AppId       string
	Fields      []string
}

func (dto CreatePermissionDto) Validate() error {
	if dto.Id == "" {
		return v1.ErrorBadRequest("empty id")
	}
	if dto.Name == "" {
		return v1.ErrorBadRequest("empty name")
	}
	if dto.GroupId == "" {
		return v1.ErrorBadRequest("empty group")
	}
	if dto.AppId == "" {
		return v1.ErrorBadRequest("empty app id")
	}
	// check if GroupId matches first part of Id
	// example: GroupId = "group1", Id = "group1.permission1"
	if dto.Id[:len(dto.GroupId)] != dto.GroupId {
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
	AppsIds    []string
	WithDenied bool
}
