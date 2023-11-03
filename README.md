# rbac

# Usage 

Service is used to manage roles and permissions.

1. Create a role
2. Create permissions
3. Assign permissions to the role
4. Assign the role to the user
5. Check if the user has permission

## Create a role
 call GRPC method `CreateRole` with role name
```protobuf
message CreateRoleRequest {
  required int64 teamID = 1;
  required string name = 2;
  required string description = 3;
}
```

## Create permissions
 call GRPC method `CreatePermission` with permission name
```protobuf
message CreatePermissionRequest {
  required string ID = 1;
  required int32 appID = 2;
  required string name = 3;
  required string description = 4;
  repeated string fields = 5;
}
```

## Assign permissions to the role
 call GRPC method `AssignPermission` with role id and permission id
```protobuf
message AssignPermissionRequest {
  required int64 roleID = 1;
  required string permissionID = 2;
}
```

## Assign the role to the user
 call GRPC method `AssignRole` with role id and user id
```protobuf
message AssignRoleRequest {
  required int64 roleID = 1;
  required int64 userID = 2;
}
```

## Check if the user has permission
 call GRPC method `CheckPermission` with user id, permission id and resource id
```protobuf
message CheckPermissionRequest {
  required int64 userID = 1;
  required string teamID = 2;
  repeated string permissionIDs = 3;
}
```


# Development
## Init project
1. ~~Change module name in go.mod~~
2. ~~Rename cmd/dummy directory~~
3. ~~Replace every "dummy" code with your code~~
4. Local env: create .env from .example, replace SERVICE_NAME, set unused HTTP_PORT, GRPC_PORT
5. CI: create .gitlab-ci.yml from .example, replace "DB_PORT: 5440X" with unused port

## Proto files

### Add a proto template

```bash
kratos proto add api/server/server.proto
```

### Generate the proto code

```bash
kratos proto client api/server/server.proto
```

### Generate the source code of service by proto file

```bash
kratos proto server api/server/server.proto -t internal/service

go generate ./...
```

## Generate other auxiliary files by Makefile

### Download and update dependencies

```bash
make init
```

### Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file

```bash
make api
```

### Generate all files

```bash
make all
```

### Generate migrations

[Install Atlas](https://entgo.io/docs/versioned-migrations#generating-migrations)

```bash
make migrations
```

## Run

### Run debug

```bash
make run
```

### Build & Run

```bash
export AWS_ACCESS_KEY_ID={aws-key}
export AWS_SECRET_ACCESS_KEY={aws-secret}

make build
```

## Run in Docker

```bash
make start
```

To stop docker:

```bash
make stop
```
