# API Proto files

## Add a proto template

```bash
kratos proto add api/rbac/v1/team-identity-role.proto
```

## Generate the proto code

```bash
kratos proto client api/rbac/v1/permissions.proto
kratos proto client api/rbac/v1/roles.proto
kratos proto client api/rbac/v1/models.proto
kratos proto client api/rbac/v1/check-permissions.proto
kratos proto client api/rbac/v1/team-identity-role.proto
kratos proto client api/rbac/v1/teams.proto
```

## Generate the source code of service by proto file

```bash
kratos proto server api/rbac/v1/permissions.proto -t internal/service
kratos proto server api/rbac/v1/roles.proto -t internal/service
kratos proto server api/rbac/v1/models.proto -t internal/service
kratos proto server api/rbac/v1/check-permissions.proto -t internal/service
kratos proto server api/rbac/v1/team-identity-role.proto -t internal/service
kratos proto server api/rbac/v1/teams.proto -t internal/service
```
