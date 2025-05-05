## Permission:

|    Field    |   Type   | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag           | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id          | string   | false  | false    | false    | false   | false         | true      | json:"id,omitempty"          |          0 |         |
| group_id    | string   | false  | false    | false    | false   | false         | true      | json:"group_id,omitempty"    |          0 |         |
| name        | string   | false  | false    | false    | false   | false         | false     | json:"name,omitempty"        |          2 |         |
| description | string   | false  | true     | false    | true    | false         | false     | json:"description,omitempty" |          0 |         |
| app_id      | string   | false  | false    | false    | false   | false         | true      | json:"app_id,omitempty"      |          1 |         |
| fields      | []string | false  | true     | false    | true    | false         | false     | json:"fields,omitempty"      |          0 |         |


| Edge  |      Type       | Inverse |   BackRef   | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| roles | RolePermission  | false   |             | O2M      | false  | true     |         |
| group | PermissionGroup | true    | permissions | M2O      | true   | false    |         |

## PermissionGroup:

| Field  |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |        StructTag        | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id     | string | false  | false    | false    | false   | false         | true      | json:"id,omitempty"     |          0 |         |
| app_id | string | false  | false    | false    | false   | false         | true      | json:"app_id,omitempty" |          1 |         |
| name   | string | false  | false    | false    | false   | false         | false     | json:"name,omitempty"   |          2 |         |


|    Edge     |    Type    | Inverse | BackRef | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| permissions | Permission | false   |         | O2M      | false  | true     |         |

## ResourceAccess:

|     Field     |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |           StructTag            | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id            | int    | false  | false    | false    | false   | false         | false     | json:"id,omitempty"            |          0 |         |
| tenant_id     | int64  | false  | false    | false    | false   | false         | true      | json:"tenant_id,omitempty"     |          0 |         |
| resource_type | string | false  | true     | true     | false   | false         | true      | json:"resource_type,omitempty" |          0 |         |
| resource_id   | int64  | false  | true     | true     | false   | false         | true      | json:"resource_id,omitempty"   |          0 |         |
| identity_id   | string | false  | false    | false    | true    | false         | true      | json:"identity_id,omitempty"   |          0 |         |
| role_id       | int64  | false  | false    | false    | false   | false         | true      | json:"role_id,omitempty"       |          0 |         |


| Edge |     Type     | Inverse | BackRef | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| role | Role         | false   |         | M2O      | true   | false    |         |
| type | ResourceType | false   |         | M2O      | true   | true     |         |

## ResourceType:

|    Field    |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag           | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id          | string | false  | false    | false    | false   | false         | true      | json:"id,omitempty"          |          0 |         |
| description | string | false  | true     | false    | true    | false         | false     | json:"description,omitempty" |          0 |         |


| Edge  | Type | Inverse | BackRef | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| roles | Role | false   |         | O2M      | false  | true     |         |

## Role:

|    Field    |   Type    | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag           | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id          | int64     | false  | false    | false    | false   | false         | true      | json:"id,omitempty"          |          0 |         |
| deleted_at  | time.Time | false  | true     | true     | false   | false         | false     | json:"deleted_at,omitempty"  |          0 |         |
| name        | string    | false  | false    | false    | false   | false         | false     | json:"name,omitempty"        |          2 |         |
| description | string    | false  | true     | false    | true    | false         | false     | json:"description,omitempty" |          0 |         |
| tenant_id   | int64     | false  | false    | false    | true    | false         | true      | json:"tenant_id,omitempty"   |          0 |         |
| is_system   | bool      | false  | false    | false    | true    | false         | true      | json:"is_system,omitempty"   |          0 |         |
| created_at  | time.Time | false  | false    | false    | true    | false         | true      | json:"created_at,omitempty"  |          0 |         |
| updated_at  | time.Time | false  | false    | false    | true    | false         | false     | json:"updated_at,omitempty"  |          0 |         |


|    Edge     |      Type      | Inverse | BackRef | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| permissions | RolePermission | false   |         | O2M      | false  | true     |         |

## RolePermission:

|     Field     |   Type   | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |           StructTag            | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id            | int      | false  | false    | false    | false   | false         | false     | json:"id,omitempty"            |          0 |         |
| tenant_id     | int64    | false  | false    | false    | true    | false         | true      | json:"tenant_id,omitempty"     |          0 |         |
| role_id       | int64    | false  | false    | false    | false   | false         | true      | json:"role_id,omitempty"       |          0 |         |
| permission_id | string   | false  | false    | false    | false   | false         | true      | json:"permission_id,omitempty" |          0 |         |
| deny          | bool     | false  | false    | false    | true    | false         | false     | json:"deny,omitempty"          |          0 |         |
| fields        | []string | false  | false    | false    | false   | false         | false     | json:"fields,omitempty"        |          0 |         |
| value         | int64    | false  | true     | false    | true    | false         | false     | json:"value,omitempty"         |          0 |         |


|    Edge    |    Type    | Inverse |   BackRef   | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| role       | Role       | true    | permissions | M2O      | true   | false    |         |
| permission | Permission | true    | roles       | M2O      | true   | false    |         |

## Team:

|    Field    |       Type        | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag           | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id          | int64             | false  | false    | false    | false   | false         | false     | json:"id,omitempty"          |          0 |         |
| deleted_at  | time.Time         | false  | true     | true     | false   | false         | false     | json:"deleted_at,omitempty"  |          0 |         |
| tenant_id   | int64             | false  | false    | false    | false   | false         | false     | json:"tenant_id,omitempty"   |          0 |         |
| parent_id   | int64             | false  | true     | true     | false   | false         | false     | json:"parent_id,omitempty"   |          0 |         |
| parents_ids | *pgtype.Int8Array | false  | true     | false    | false   | false         | false     | json:"parents_ids,omitempty" |          0 |         |
| name        | string            | false  | false    | false    | false   | false         | false     | json:"name,omitempty"        |          0 |         |
| description | string            | false  | true     | false    | true    | false         | false     | json:"description,omitempty" |          0 |         |
| created_at  | time.Time         | false  | false    | false    | true    | false         | false     | json:"created_at,omitempty"  |          0 |         |
| updated_at  | time.Time         | false  | false    | false    | true    | false         | false     | json:"updated_at,omitempty"  |          0 |         |


|   Edge   | Type | Inverse | BackRef  | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| parent   | Team | true    | children | M2O      | true   | true     |         |
| children | Team | false   |          | O2M      | false  | true     |         |

## TeamIdentityRole:

|    Field    |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag           | Validators | Comment |
|---|---|---|---|---|---|---|---|---|---|---|
| id          | int    | false  | false    | false    | false   | false         | false     | json:"id,omitempty"          |          0 |         |
| tenant_id   | int64  | false  | false    | false    | false   | false         | true      | json:"tenant_id,omitempty"   |          0 |         |
| team_id     | int64  | false  | true     | true     | false   | false         | true      | json:"team_id,omitempty"     |          0 |         |
| identity_id | string | false  | false    | false    | true    | false         | true      | json:"identity_id,omitempty" |          0 |         |
| role_id     | int64  | false  | false    | false    | false   | false         | true      | json:"role_id,omitempty"     |          0 |         |


| Edge | Type | Inverse | BackRef | Relation | Unique | Optional | Comment |
|---|---|---|---|---|---|---|---|
| role | Role | false   |         | M2O      | true   | false    |         |
| team | Team | false   |         | M2O      | true   | true     |         |

