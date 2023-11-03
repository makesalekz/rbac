
### Permissions

| Col         | Description                                   | type        |
|-------------|-----------------------------------------------|-------------|
| id          | Permission ID in the format app:entity:action | *string     |
| app_id      | Application string - ID                       | *string     |
| name        | Readable name                                 | *string(32) |
| description | Readable description                          | string      |
| fields      | All entity fields                             | []string    |

Permissions for each application in which access to entities and fields is described for each application, entity and action
Additionally described all entity fields for validation

### Roles
| Col         | Description                      | type        |
|-------------|----------------------------------|-------------|
| id          | Role ID                          | *int        |
| tenant_id   | tenant ID                        | *int        |
| name        | Readable name                    | *string(32) |
| description | Readable description             | string      |
| is_system   | Is system role (default = false) | bool        |

Roles for tenants
Roles can be system - tenant_id is null

### RolePermissions
| Col           | Description           | type     |
|---------------|-----------------------|----------|
| tenant_id     | Tenant ID             | int      |
| role_id       | Role ID               | *int     |
| permission_id | Permission ID         | *string  |
| deny          | bool                  | bool     |
| fields        | Fields for permission | []string |

Mapping between roles and permissions
Column deny - deny access to field, if deny = true, access to field is denied

### Teams
| Col        | Description                         | type    |
|------------|-------------------------------------|---------|
| id         | Team ID                             | *int    |
| tenant_id  | ID of tenant                        | *int    |
| parent_id  | the closest parent to this command. | int     |
| parent_ids | all parent teams                    | int     |
| name       | Readable name of team               | *string |

Tenant teams with hierarchy

### TeamIdentityRole

| Col         | Description           | type |
|-------------|-----------------------|------|
| tenant_id   | Tenant ID             | *int |
| team_id     | Team ID               | int  |
| identity_id | Identity ID           | int  |
| role_id     | Role ID               | *int |


Mapping between teams, identities and roles
team_id and identity_id can be null it means that role is default for tenant and all identities, every user has default role for tenant

      * - can not be null (required)

### How to use?

