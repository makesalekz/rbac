-- Create basic roles
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at") VALUES
  (1, 'Admin', 'Has all permissions in a tenant', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (2, 'Basic', 'Has all limitations according to the Basic plan', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Set sequence to 1000 to avoid conflicts with system roles
ALTER SEQUENCE "roles_id_seq" RESTART WITH 1000;

-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields") VALUES
-- Tenant admin role
  (1, 'admin.permission.read', FALSE, '[]'),
  (1, 'admin.role.create', FALSE, '[]'),
  (1, 'admin.role.read', FALSE, '[]'),
  (1, 'admin.role.update', FALSE, '[]'),
  (1, 'admin.role.delete', FALSE, '[]'),
  (1, 'admin.team.create', FALSE, '[]'),
  (1, 'admin.team.read', FALSE, '[]'),
  (1, 'admin.team.update', FALSE, '[]'),
  (1, 'admin.team.delete', FALSE, '[]'),
  (1, 'admin.role.assign', FALSE, '[]'),
  (1, 'admin.tenant.update', FALSE, '[]'),
  (1, 'admin.tenant.delete', FALSE, '[]'),
  (1, 'admin.group.create', FALSE, '[]'),
  (1, 'admin.group.read', FALSE, '[]'),
  (1, 'admin.group.update', FALSE, '[]'),
  (1, 'admin.group.delete', FALSE, '[]'),
  (1, 'admin.member.read', FALSE, '[]'),
  (1, 'admin.member.delete', FALSE, '[]'),
  (1, 'admin.invite.create', FALSE, '[]'),
  (1, 'admin.invite.read', FALSE, '[]'),
  (1, 'admin.invite.delete', FALSE, '[]'),

-- Basic plan limitations
  (2, 'admin.permission.create', TRUE, '[]'),
  (2, 'admin.permission.update', TRUE, '[]'),
  (2, 'admin.permission.delete', TRUE, '[]'),
  (2, 'admin.role_system.create', TRUE, '[]'),
  (2, 'admin.role_system.update', TRUE, '[]'),
  (2, 'admin.role_system.delete', TRUE, '[]'),
  (2, 'admin.team.create', TRUE, '[]'),
  (2, 'admin.team.read', TRUE, '[]'),
  (2, 'admin.team.update', TRUE, '[]'),
  (2, 'admin.team.delete', TRUE, '[]');
