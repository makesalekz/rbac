<!-- Generator: Widdershins v4.0.1 -->

<h1 id="api"> v0.0.1</h1>

> Scroll down for example requests and responses.

<h1 id="api-assigns">Assigns</h1>

## Assigns_AssignRole

<a id="opIdAssigns_AssignRole"></a>

`POST /v1/rbac/assigns`

> Body parameter

```json
{
  "identityId": "string",
  "roleId": "string",
  "teamId": "string"
}
```

<h3 id="assigns_assignrole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.AssignRoleRequest](#schemarbac.v1.assignrolerequest)|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="assigns_assignrole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Assigns_ListAssigns

<a id="opIdAssigns_ListAssigns"></a>

`POST /v1/rbac/assigns/list`

> Body parameter

```json
{
  "identityId": "string",
  "teamId": "string"
}
```

<h3 id="assigns_listassigns-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.ListAssignsRequest](#schemarbac.v1.listassignsrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "roles": [
    {
      "assignId": "string",
      "role": {
        "id": "string",
        "name": "string",
        "description": "string",
        "isSystem": true,
        "createdAt": "string",
        "updatedAt": "string",
        "deletedAt": "string"
      },
      "identityId": "string",
      "team": {
        "id": "string",
        "ownerId": "string",
        "parentId": "string",
        "parentsIds": [
          "string"
        ],
        "name": "string",
        "description": "string",
        "createdAt": "string",
        "updatedAt": "string",
        "subs": [
          {}
        ]
      }
    }
  ]
}
```

<h3 id="assigns_listassigns-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.ListAssignsReply](#schemarbac.v1.listassignsreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Assigns_UnassignRole

<a id="opIdAssigns_UnassignRole"></a>

`DELETE /v1/rbac/assigns/{assignId}`

<h3 id="assigns_unassignrole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|assignId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="assigns_unassignrole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="api-checkpermissions">CheckPermissions</h1>

## CheckPermissions_CheckPermissions

<a id="opIdCheckPermissions_CheckPermissions"></a>

`POST /v1/rbac/check`

> Body parameter

```json
{
  "tenantId": "string",
  "teamId": "string",
  "permissions": [
    "string"
  ],
  "identities": [
    "string"
  ]
}
```

<h3 id="checkpermissions_checkpermissions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.CheckPermissionsRequest](#schemarbac.v1.checkpermissionsrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "permissions": {
    "property1": {
      "fields": [
        "string"
      ]
    },
    "property2": {
      "fields": [
        "string"
      ]
    }
  }
}
```

<h3 id="checkpermissions_checkpermissions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.CheckPermissionsReply](#schemarbac.v1.checkpermissionsreply)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="api-permissions">Permissions</h1>

## Permissions_CreatePermission

<a id="opIdPermissions_CreatePermission"></a>

`POST /v1/rbac/permissions`

> Body parameter

```json
{
  "id": "string",
  "groupId": "string",
  "appId": "string",
  "name": "string",
  "description": "string",
  "fields": [
    "string"
  ]
}
```

<h3 id="permissions_createpermission-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.CreatePermissionRequest](#schemarbac.v1.createpermissionrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "permission": {
    "id": "string",
    "name": "string",
    "description": "string",
    "appId": "string",
    "groupId": "string",
    "fields": [
      "string"
    ]
  }
}
```

<h3 id="permissions_createpermission-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.PermissionReply](#schemarbac.v1.permissionreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Permissions_ListPermissions

<a id="opIdPermissions_ListPermissions"></a>

`POST /v1/rbac/permissions/list`

> Body parameter

```json
{
  "appsIds": [
    "string"
  ]
}
```

<h3 id="permissions_listpermissions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.ListPermissionsRequest](#schemarbac.v1.listpermissionsrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "groups": [
    {
      "id": "string",
      "name": "string",
      "appId": "string",
      "permissions": [
        {
          "id": "string",
          "name": "string",
          "description": "string",
          "appId": "string",
          "groupId": "string",
          "fields": [
            "string"
          ]
        }
      ]
    }
  ]
}
```

<h3 id="permissions_listpermissions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.ListPermissionsReply](#schemarbac.v1.listpermissionsreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Permissions_GetPermission

<a id="opIdPermissions_GetPermission"></a>

`GET /v1/rbac/permissions/{permissionId}`

<h3 id="permissions_getpermission-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|permissionId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "permission": {
    "id": "string",
    "name": "string",
    "description": "string",
    "appId": "string",
    "groupId": "string",
    "fields": [
      "string"
    ]
  }
}
```

<h3 id="permissions_getpermission-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.PermissionReply](#schemarbac.v1.permissionreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Permissions_UpdatePermission

<a id="opIdPermissions_UpdatePermission"></a>

`PUT /v1/rbac/permissions/{permissionId}`

> Body parameter

```json
{
  "permissionId": "string",
  "name": "string",
  "description": "string",
  "fields": [
    "string"
  ]
}
```

<h3 id="permissions_updatepermission-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|permissionId|path|string|true|none|
|body|body|[rbac.v1.UpdatePermissionRequest](#schemarbac.v1.updatepermissionrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "permission": {
    "id": "string",
    "name": "string",
    "description": "string",
    "appId": "string",
    "groupId": "string",
    "fields": [
      "string"
    ]
  }
}
```

<h3 id="permissions_updatepermission-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.PermissionReply](#schemarbac.v1.permissionreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Permissions_DeletePermission

<a id="opIdPermissions_DeletePermission"></a>

`DELETE /v1/rbac/permissions/{permissionId}`

<h3 id="permissions_deletepermission-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|permissionId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="permissions_deletepermission-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="api-roles">Roles</h1>

## Roles_CreateRole

<a id="opIdRoles_CreateRole"></a>

`POST /v1/rbac/roles`

> Body parameter

```json
{
  "name": "string",
  "description": "string",
  "isSystem": true,
  "allow": [
    "string"
  ],
  "deny": [
    "string"
  ]
}
```

<h3 id="roles_createrole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.CreateRoleRequest](#schemarbac.v1.createrolerequest)|true|none|

> Example responses

> 200 Response

```json
{
  "role": {
    "id": "string",
    "name": "string",
    "description": "string",
    "isSystem": true,
    "createdAt": "string",
    "updatedAt": "string",
    "deletedAt": "string"
  }
}
```

<h3 id="roles_createrole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.RoleReply](#schemarbac.v1.rolereply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_ListRoles

<a id="opIdRoles_ListRoles"></a>

`POST /v1/rbac/roles/list`

> Body parameter

```json
{
  "search": "string"
}
```

<h3 id="roles_listroles-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.ListRolesRequest](#schemarbac.v1.listrolesrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "roles": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "isSystem": true,
      "createdAt": "string",
      "updatedAt": "string",
      "deletedAt": "string"
    }
  ]
}
```

<h3 id="roles_listroles-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.ListRolesReply](#schemarbac.v1.listrolesreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_GetRole

<a id="opIdRoles_GetRole"></a>

`GET /v1/rbac/roles/{roleId}`

<h3 id="roles_getrole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "role": {
    "id": "string",
    "name": "string",
    "description": "string",
    "isSystem": true,
    "createdAt": "string",
    "updatedAt": "string",
    "deletedAt": "string"
  }
}
```

<h3 id="roles_getrole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.RoleReply](#schemarbac.v1.rolereply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_UpdateRole

<a id="opIdRoles_UpdateRole"></a>

`PUT /v1/rbac/roles/{roleId}`

> Body parameter

```json
{
  "roleId": "string",
  "name": "string",
  "description": "string",
  "allow": [
    "string"
  ],
  "deny": [
    "string"
  ]
}
```

<h3 id="roles_updaterole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|
|body|body|[rbac.v1.UpdateRoleRequest](#schemarbac.v1.updaterolerequest)|true|none|

> Example responses

> 200 Response

```json
{
  "role": {
    "id": "string",
    "name": "string",
    "description": "string",
    "isSystem": true,
    "createdAt": "string",
    "updatedAt": "string",
    "deletedAt": "string"
  }
}
```

<h3 id="roles_updaterole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.RoleReply](#schemarbac.v1.rolereply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_DeleteRole

<a id="opIdRoles_DeleteRole"></a>

`DELETE /v1/rbac/roles/{roleId}`

<h3 id="roles_deleterole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="roles_deleterole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_AddPermissionToRole

<a id="opIdRoles_AddPermissionToRole"></a>

`POST /v1/rbac/roles/{roleId}/permissions`

> Body parameter

```json
{
  "roleId": "string",
  "permissionId": "string",
  "deny": true,
  "fields": [
    "string"
  ]
}
```

<h3 id="roles_addpermissiontorole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|
|body|body|[rbac.v1.AddPermissionToRoleRequest](#schemarbac.v1.addpermissiontorolerequest)|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="roles_addpermissiontorole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_ListRolePermissions

<a id="opIdRoles_ListRolePermissions"></a>

`POST /v1/rbac/roles/{roleId}/permissions/list`

> Body parameter

```json
{
  "roleId": "string"
}
```

<h3 id="roles_listrolepermissions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|
|body|body|[rbac.v1.RoleRequest](#schemarbac.v1.rolerequest)|true|none|

> Example responses

> 200 Response

```json
{
  "permissions": [
    {
      "id": "string",
      "deny": true,
      "fields": [
        "string"
      ]
    }
  ]
}
```

<h3 id="roles_listrolepermissions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.RolePermissionsReply](#schemarbac.v1.rolepermissionsreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Roles_RemovePermissionFromRole

<a id="opIdRoles_RemovePermissionFromRole"></a>

`DELETE /v1/rbac/roles/{roleId}/permissions/{permissionId}`

<h3 id="roles_removepermissionfromrole-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|roleId|path|string|true|none|
|permissionId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="roles_removepermissionfromrole-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="api-teams">Teams</h1>

## Teams_CreateTeam

<a id="opIdTeams_CreateTeam"></a>

`POST /v1/rbac/teams`

> Body parameter

```json
{
  "name": "string",
  "description": "string",
  "parentId": "string"
}
```

<h3 id="teams_createteam-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.CreateTeamRequest](#schemarbac.v1.createteamrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "team": {
    "id": "string",
    "ownerId": "string",
    "parentId": "string",
    "parentsIds": [
      "string"
    ],
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "subs": [
      {}
    ]
  },
  "memberId": "string"
}
```

<h3 id="teams_createteam-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.TeamReply](#schemarbac.v1.teamreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Teams_ListTeams

<a id="opIdTeams_ListTeams"></a>

`POST /v1/rbac/teams/list`

> Body parameter

```json
{
  "parentId": "string",
  "paginate": {
    "limit": 0,
    "fromId": "string",
    "toId": "string",
    "aroundId": "string",
    "fromDate": "string",
    "toDate": "string",
    "aroundDate": "string",
    "fromLabel": "string",
    "descending": true,
    "page": 0
  }
}
```

<h3 id="teams_listteams-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[rbac.v1.ListTeamsRequest](#schemarbac.v1.listteamsrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "teams": [
    {
      "id": "string",
      "ownerId": "string",
      "parentId": "string",
      "parentsIds": [
        "string"
      ],
      "name": "string",
      "description": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "subs": [
        {}
      ]
    }
  ],
  "paginate": {
    "total": 0,
    "fromId": "string",
    "toId": "string",
    "fromDate": "string",
    "toDate": "string",
    "fromLabel": "string"
  }
}
```

<h3 id="teams_listteams-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.ListTeamsReply](#schemarbac.v1.listteamsreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Teams_GetTeam

<a id="opIdTeams_GetTeam"></a>

`GET /v1/rbac/teams/{teamId}`

<h3 id="teams_getteam-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|teamId|path|string|true|none|
|withTree|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "team": {
    "id": "string",
    "ownerId": "string",
    "parentId": "string",
    "parentsIds": [
      "string"
    ],
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "subs": [
      {}
    ]
  },
  "memberId": "string"
}
```

<h3 id="teams_getteam-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.TeamReply](#schemarbac.v1.teamreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Teams_UpdateTeam

<a id="opIdTeams_UpdateTeam"></a>

`PUT /v1/rbac/teams/{teamId}`

> Body parameter

```json
{
  "teamId": "string",
  "name": "string",
  "description": "string"
}
```

<h3 id="teams_updateteam-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|teamId|path|string|true|none|
|body|body|[rbac.v1.UpdateTeamRequest](#schemarbac.v1.updateteamrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "team": {
    "id": "string",
    "ownerId": "string",
    "parentId": "string",
    "parentsIds": [
      "string"
    ],
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "subs": [
      {}
    ]
  },
  "memberId": "string"
}
```

<h3 id="teams_updateteam-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[rbac.v1.TeamReply](#schemarbac.v1.teamreply)|

<aside class="success">
This operation does not require authentication
</aside>

## Teams_DeleteTeam

<a id="opIdTeams_DeleteTeam"></a>

`DELETE /v1/rbac/teams/{teamId}`

<h3 id="teams_deleteteam-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|teamId|path|string|true|none|
|withTree|query|boolean|false|none|

> Example responses

> 200 Response

```json
{}
```

<h3 id="teams_deleteteam-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[utils.v1.EmptyReply](#schemautils.v1.emptyreply)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_rbac.v1.AddPermissionToRoleRequest">rbac.v1.AddPermissionToRoleRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.addpermissiontorolerequest"></a>
<a id="schema_rbac.v1.AddPermissionToRoleRequest"></a>
<a id="tocSrbac.v1.addpermissiontorolerequest"></a>
<a id="tocsrbac.v1.addpermissiontorolerequest"></a>

```json
{
  "roleId": "string",
  "permissionId": "string",
  "deny": true,
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|roleId|string|false|none|none|
|permissionId|string|false|none|none|
|deny|boolean|false|none|none|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.AssignRoleRequest">rbac.v1.AssignRoleRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.assignrolerequest"></a>
<a id="schema_rbac.v1.AssignRoleRequest"></a>
<a id="tocSrbac.v1.assignrolerequest"></a>
<a id="tocsrbac.v1.assignrolerequest"></a>

```json
{
  "identityId": "string",
  "roleId": "string",
  "teamId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|identityId|string|false|none|none|
|roleId|string|false|none|none|
|teamId|string|false|none|none|

<h2 id="tocS_rbac.v1.AssignedRole">rbac.v1.AssignedRole</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.assignedrole"></a>
<a id="schema_rbac.v1.AssignedRole"></a>
<a id="tocSrbac.v1.assignedrole"></a>
<a id="tocsrbac.v1.assignedrole"></a>

```json
{
  "assignId": "string",
  "role": {
    "id": "string",
    "name": "string",
    "description": "string",
    "isSystem": true,
    "createdAt": "string",
    "updatedAt": "string",
    "deletedAt": "string"
  },
  "identityId": "string",
  "team": {
    "id": "string",
    "ownerId": "string",
    "parentId": "string",
    "parentsIds": [
      "string"
    ],
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "subs": [
      {}
    ]
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|assignId|string|false|none|none|
|role|[rbac.v1.Role](#schemarbac.v1.role)|false|none|none|
|identityId|string|false|none|none|
|team|[rbac.v1.Team](#schemarbac.v1.team)|false|none|none|

<h2 id="tocS_rbac.v1.CheckPermissionsReply">rbac.v1.CheckPermissionsReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.checkpermissionsreply"></a>
<a id="schema_rbac.v1.CheckPermissionsReply"></a>
<a id="tocSrbac.v1.checkpermissionsreply"></a>
<a id="tocsrbac.v1.checkpermissionsreply"></a>

```json
{
  "permissions": {
    "property1": {
      "fields": [
        "string"
      ]
    },
    "property2": {
      "fields": [
        "string"
      ]
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|permissions|object|false|none|none|
|» **additionalProperties**|[rbac.v1.ListOfFields](#schemarbac.v1.listoffields)|false|none|none|

<h2 id="tocS_rbac.v1.CheckPermissionsRequest">rbac.v1.CheckPermissionsRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.checkpermissionsrequest"></a>
<a id="schema_rbac.v1.CheckPermissionsRequest"></a>
<a id="tocSrbac.v1.checkpermissionsrequest"></a>
<a id="tocsrbac.v1.checkpermissionsrequest"></a>

```json
{
  "tenantId": "string",
  "teamId": "string",
  "permissions": [
    "string"
  ],
  "identities": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|tenantId|string|false|none|none|
|teamId|string|false|none|none|
|permissions|[string]|false|none|none|
|identities|[string]|false|none|none|

<h2 id="tocS_rbac.v1.CreatePermissionRequest">rbac.v1.CreatePermissionRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.createpermissionrequest"></a>
<a id="schema_rbac.v1.CreatePermissionRequest"></a>
<a id="tocSrbac.v1.createpermissionrequest"></a>
<a id="tocsrbac.v1.createpermissionrequest"></a>

```json
{
  "id": "string",
  "groupId": "string",
  "appId": "string",
  "name": "string",
  "description": "string",
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|groupId|string|false|none|none|
|appId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.CreateRoleRequest">rbac.v1.CreateRoleRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.createrolerequest"></a>
<a id="schema_rbac.v1.CreateRoleRequest"></a>
<a id="tocSrbac.v1.createrolerequest"></a>
<a id="tocsrbac.v1.createrolerequest"></a>

```json
{
  "name": "string",
  "description": "string",
  "isSystem": true,
  "allow": [
    "string"
  ],
  "deny": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|none|
|description|string|false|none|none|
|isSystem|boolean|false|none|none|
|allow|[string]|false|none|none|
|deny|[string]|false|none|none|

<h2 id="tocS_rbac.v1.CreateTeamRequest">rbac.v1.CreateTeamRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.createteamrequest"></a>
<a id="schema_rbac.v1.CreateTeamRequest"></a>
<a id="tocSrbac.v1.createteamrequest"></a>
<a id="tocsrbac.v1.createteamrequest"></a>

```json
{
  "name": "string",
  "description": "string",
  "parentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|none|
|description|string|false|none|none|
|parentId|string|false|none|none|

<h2 id="tocS_rbac.v1.Group">rbac.v1.Group</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.group"></a>
<a id="schema_rbac.v1.Group"></a>
<a id="tocSrbac.v1.group"></a>
<a id="tocsrbac.v1.group"></a>

```json
{
  "id": "string",
  "name": "string",
  "appId": "string",
  "permissions": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "appId": "string",
      "groupId": "string",
      "fields": [
        "string"
      ]
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|name|string|false|none|none|
|appId|string|false|none|none|
|permissions|[[rbac.v1.Permission](#schemarbac.v1.permission)]|false|none|none|

<h2 id="tocS_rbac.v1.ListAssignsReply">rbac.v1.ListAssignsReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listassignsreply"></a>
<a id="schema_rbac.v1.ListAssignsReply"></a>
<a id="tocSrbac.v1.listassignsreply"></a>
<a id="tocsrbac.v1.listassignsreply"></a>

```json
{
  "roles": [
    {
      "assignId": "string",
      "role": {
        "id": "string",
        "name": "string",
        "description": "string",
        "isSystem": true,
        "createdAt": "string",
        "updatedAt": "string",
        "deletedAt": "string"
      },
      "identityId": "string",
      "team": {
        "id": "string",
        "ownerId": "string",
        "parentId": "string",
        "parentsIds": [
          "string"
        ],
        "name": "string",
        "description": "string",
        "createdAt": "string",
        "updatedAt": "string",
        "subs": [
          {}
        ]
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|roles|[[rbac.v1.AssignedRole](#schemarbac.v1.assignedrole)]|false|none|none|

<h2 id="tocS_rbac.v1.ListAssignsRequest">rbac.v1.ListAssignsRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listassignsrequest"></a>
<a id="schema_rbac.v1.ListAssignsRequest"></a>
<a id="tocSrbac.v1.listassignsrequest"></a>
<a id="tocsrbac.v1.listassignsrequest"></a>

```json
{
  "identityId": "string",
  "teamId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|identityId|string|false|none|none|
|teamId|string|false|none|none|

<h2 id="tocS_rbac.v1.ListOfFields">rbac.v1.ListOfFields</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listoffields"></a>
<a id="schema_rbac.v1.ListOfFields"></a>
<a id="tocSrbac.v1.listoffields"></a>
<a id="tocsrbac.v1.listoffields"></a>

```json
{
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.ListPermissionsReply">rbac.v1.ListPermissionsReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listpermissionsreply"></a>
<a id="schema_rbac.v1.ListPermissionsReply"></a>
<a id="tocSrbac.v1.listpermissionsreply"></a>
<a id="tocsrbac.v1.listpermissionsreply"></a>

```json
{
  "groups": [
    {
      "id": "string",
      "name": "string",
      "appId": "string",
      "permissions": [
        {
          "id": "string",
          "name": "string",
          "description": "string",
          "appId": "string",
          "groupId": "string",
          "fields": [
            "string"
          ]
        }
      ]
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|groups|[[rbac.v1.Group](#schemarbac.v1.group)]|false|none|none|

<h2 id="tocS_rbac.v1.ListPermissionsRequest">rbac.v1.ListPermissionsRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listpermissionsrequest"></a>
<a id="schema_rbac.v1.ListPermissionsRequest"></a>
<a id="tocSrbac.v1.listpermissionsrequest"></a>
<a id="tocsrbac.v1.listpermissionsrequest"></a>

```json
{
  "appsIds": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|appsIds|[string]|false|none|none|

<h2 id="tocS_rbac.v1.ListRolesReply">rbac.v1.ListRolesReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listrolesreply"></a>
<a id="schema_rbac.v1.ListRolesReply"></a>
<a id="tocSrbac.v1.listrolesreply"></a>
<a id="tocsrbac.v1.listrolesreply"></a>

```json
{
  "roles": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "isSystem": true,
      "createdAt": "string",
      "updatedAt": "string",
      "deletedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|roles|[[rbac.v1.Role](#schemarbac.v1.role)]|false|none|none|

<h2 id="tocS_rbac.v1.ListRolesRequest">rbac.v1.ListRolesRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listrolesrequest"></a>
<a id="schema_rbac.v1.ListRolesRequest"></a>
<a id="tocSrbac.v1.listrolesrequest"></a>
<a id="tocsrbac.v1.listrolesrequest"></a>

```json
{
  "search": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|search|string|false|none|none|

<h2 id="tocS_rbac.v1.ListTeamsReply">rbac.v1.ListTeamsReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listteamsreply"></a>
<a id="schema_rbac.v1.ListTeamsReply"></a>
<a id="tocSrbac.v1.listteamsreply"></a>
<a id="tocsrbac.v1.listteamsreply"></a>

```json
{
  "teams": [
    {
      "id": "string",
      "ownerId": "string",
      "parentId": "string",
      "parentsIds": [
        "string"
      ],
      "name": "string",
      "description": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "subs": [
        {}
      ]
    }
  ],
  "paginate": {
    "total": 0,
    "fromId": "string",
    "toId": "string",
    "fromDate": "string",
    "toDate": "string",
    "fromLabel": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|teams|[[rbac.v1.Team](#schemarbac.v1.team)]|false|none|none|
|paginate|[utils.v1.PaginateReply](#schemautils.v1.paginatereply)|false|none|none|

<h2 id="tocS_rbac.v1.ListTeamsRequest">rbac.v1.ListTeamsRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.listteamsrequest"></a>
<a id="schema_rbac.v1.ListTeamsRequest"></a>
<a id="tocSrbac.v1.listteamsrequest"></a>
<a id="tocsrbac.v1.listteamsrequest"></a>

```json
{
  "parentId": "string",
  "paginate": {
    "limit": 0,
    "fromId": "string",
    "toId": "string",
    "aroundId": "string",
    "fromDate": "string",
    "toDate": "string",
    "aroundDate": "string",
    "fromLabel": "string",
    "descending": true,
    "page": 0
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|parentId|string|false|none|none|
|paginate|[utils.v1.PaginateRequest](#schemautils.v1.paginaterequest)|false|none|none|

<h2 id="tocS_rbac.v1.Permission">rbac.v1.Permission</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.permission"></a>
<a id="schema_rbac.v1.Permission"></a>
<a id="tocSrbac.v1.permission"></a>
<a id="tocsrbac.v1.permission"></a>

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "appId": "string",
  "groupId": "string",
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|appId|string|false|none|none|
|groupId|string|false|none|none|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.PermissionReply">rbac.v1.PermissionReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.permissionreply"></a>
<a id="schema_rbac.v1.PermissionReply"></a>
<a id="tocSrbac.v1.permissionreply"></a>
<a id="tocsrbac.v1.permissionreply"></a>

```json
{
  "permission": {
    "id": "string",
    "name": "string",
    "description": "string",
    "appId": "string",
    "groupId": "string",
    "fields": [
      "string"
    ]
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|permission|[rbac.v1.Permission](#schemarbac.v1.permission)|false|none|none|

<h2 id="tocS_rbac.v1.Role">rbac.v1.Role</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.role"></a>
<a id="schema_rbac.v1.Role"></a>
<a id="tocSrbac.v1.role"></a>
<a id="tocsrbac.v1.role"></a>

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "isSystem": true,
  "createdAt": "string",
  "updatedAt": "string",
  "deletedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|isSystem|boolean|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|deletedAt|string|false|none|none|

<h2 id="tocS_rbac.v1.RolePermission">rbac.v1.RolePermission</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.rolepermission"></a>
<a id="schema_rbac.v1.RolePermission"></a>
<a id="tocSrbac.v1.rolepermission"></a>
<a id="tocsrbac.v1.rolepermission"></a>

```json
{
  "id": "string",
  "deny": true,
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|deny|boolean|false|none|none|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.RolePermissionsReply">rbac.v1.RolePermissionsReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.rolepermissionsreply"></a>
<a id="schema_rbac.v1.RolePermissionsReply"></a>
<a id="tocSrbac.v1.rolepermissionsreply"></a>
<a id="tocsrbac.v1.rolepermissionsreply"></a>

```json
{
  "permissions": [
    {
      "id": "string",
      "deny": true,
      "fields": [
        "string"
      ]
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|permissions|[[rbac.v1.RolePermission](#schemarbac.v1.rolepermission)]|false|none|none|

<h2 id="tocS_rbac.v1.RoleReply">rbac.v1.RoleReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.rolereply"></a>
<a id="schema_rbac.v1.RoleReply"></a>
<a id="tocSrbac.v1.rolereply"></a>
<a id="tocsrbac.v1.rolereply"></a>

```json
{
  "role": {
    "id": "string",
    "name": "string",
    "description": "string",
    "isSystem": true,
    "createdAt": "string",
    "updatedAt": "string",
    "deletedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|role|[rbac.v1.Role](#schemarbac.v1.role)|false|none|none|

<h2 id="tocS_rbac.v1.RoleRequest">rbac.v1.RoleRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.rolerequest"></a>
<a id="schema_rbac.v1.RoleRequest"></a>
<a id="tocSrbac.v1.rolerequest"></a>
<a id="tocsrbac.v1.rolerequest"></a>

```json
{
  "roleId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|roleId|string|false|none|none|

<h2 id="tocS_rbac.v1.Team">rbac.v1.Team</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.team"></a>
<a id="schema_rbac.v1.Team"></a>
<a id="tocSrbac.v1.team"></a>
<a id="tocsrbac.v1.team"></a>

```json
{
  "id": "string",
  "ownerId": "string",
  "parentId": "string",
  "parentsIds": [
    "string"
  ],
  "name": "string",
  "description": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "subs": [
    {
      "id": "string",
      "ownerId": "string",
      "parentId": "string",
      "parentsIds": [
        "string"
      ],
      "name": "string",
      "description": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "subs": []
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|ownerId|string|false|none|none|
|parentId|string|false|none|none|
|parentsIds|[string]|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|subs|[[rbac.v1.Team](#schemarbac.v1.team)]|false|none|none|

<h2 id="tocS_rbac.v1.TeamReply">rbac.v1.TeamReply</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.teamreply"></a>
<a id="schema_rbac.v1.TeamReply"></a>
<a id="tocSrbac.v1.teamreply"></a>
<a id="tocsrbac.v1.teamreply"></a>

```json
{
  "team": {
    "id": "string",
    "ownerId": "string",
    "parentId": "string",
    "parentsIds": [
      "string"
    ],
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "subs": [
      {}
    ]
  },
  "memberId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|team|[rbac.v1.Team](#schemarbac.v1.team)|false|none|none|
|memberId|string|false|none|none|

<h2 id="tocS_rbac.v1.UpdatePermissionRequest">rbac.v1.UpdatePermissionRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.updatepermissionrequest"></a>
<a id="schema_rbac.v1.UpdatePermissionRequest"></a>
<a id="tocSrbac.v1.updatepermissionrequest"></a>
<a id="tocsrbac.v1.updatepermissionrequest"></a>

```json
{
  "permissionId": "string",
  "name": "string",
  "description": "string",
  "fields": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|permissionId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|fields|[string]|false|none|none|

<h2 id="tocS_rbac.v1.UpdateRoleRequest">rbac.v1.UpdateRoleRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.updaterolerequest"></a>
<a id="schema_rbac.v1.UpdateRoleRequest"></a>
<a id="tocSrbac.v1.updaterolerequest"></a>
<a id="tocsrbac.v1.updaterolerequest"></a>

```json
{
  "roleId": "string",
  "name": "string",
  "description": "string",
  "allow": [
    "string"
  ],
  "deny": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|roleId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|allow|[string]|false|none|none|
|deny|[string]|false|none|none|

<h2 id="tocS_rbac.v1.UpdateTeamRequest">rbac.v1.UpdateTeamRequest</h2>
<!-- backwards compatibility -->
<a id="schemarbac.v1.updateteamrequest"></a>
<a id="schema_rbac.v1.UpdateTeamRequest"></a>
<a id="tocSrbac.v1.updateteamrequest"></a>
<a id="tocsrbac.v1.updateteamrequest"></a>

```json
{
  "teamId": "string",
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|teamId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|

<h2 id="tocS_utils.v1.EmptyReply">utils.v1.EmptyReply</h2>
<!-- backwards compatibility -->
<a id="schemautils.v1.emptyreply"></a>
<a id="schema_utils.v1.EmptyReply"></a>
<a id="tocSutils.v1.emptyreply"></a>
<a id="tocsutils.v1.emptyreply"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_utils.v1.PaginateReply">utils.v1.PaginateReply</h2>
<!-- backwards compatibility -->
<a id="schemautils.v1.paginatereply"></a>
<a id="schema_utils.v1.PaginateReply"></a>
<a id="tocSutils.v1.paginatereply"></a>
<a id="tocsutils.v1.paginatereply"></a>

```json
{
  "total": 0,
  "fromId": "string",
  "toId": "string",
  "fromDate": "string",
  "toDate": "string",
  "fromLabel": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|integer(int32)|false|none|none|
|fromId|string|false|none|none|
|toId|string|false|none|none|
|fromDate|string|false|none|none|
|toDate|string|false|none|none|
|fromLabel|string|false|none|none|

<h2 id="tocS_utils.v1.PaginateRequest">utils.v1.PaginateRequest</h2>
<!-- backwards compatibility -->
<a id="schemautils.v1.paginaterequest"></a>
<a id="schema_utils.v1.PaginateRequest"></a>
<a id="tocSutils.v1.paginaterequest"></a>
<a id="tocsutils.v1.paginaterequest"></a>

```json
{
  "limit": 0,
  "fromId": "string",
  "toId": "string",
  "aroundId": "string",
  "fromDate": "string",
  "toDate": "string",
  "aroundDate": "string",
  "fromLabel": "string",
  "descending": true,
  "page": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|limit|integer(int32)|false|none|none|
|fromId|string|false|none|none|
|toId|string|false|none|none|
|aroundId|string|false|none|none|
|fromDate|string|false|none|none|
|toDate|string|false|none|none|
|aroundDate|string|false|none|none|
|fromLabel|string|false|none|none|
|descending|boolean|false|none|none|
|page|integer(int32)|false|none|none|

