# API Proto files

## Add a proto template

```bash
kratos proto add api/users-roles/v1/users-roles.proto
```

## Generate the proto code

```bash
kratos proto client api/permissions/v1/permissions.proto
kratos proto client api/roles/v1/roles.proto
kratos proto client api/users-roles/v1/users-roles.proto
kratos proto client api/check-permissions/v1/check-permissions.proto
```

## Generate the source code of service by proto file

```bash
kratos proto server api/server/v1/server.proto -t internal/service
```
