-- Create resource types
INSERT INTO "resource_types" ("id", "description") VALUES
  ('team', 'RBAC: Team'),
  ('project', 'Projects: Project'),
  ('task', 'Tasks: Task');

-- Update resource permissions
INSERT INTO "resource_accesses" ("tenant_id", "resource_type", "resource_id", "identity_id", "role_id")
  SELECT "tenant_id", 'team', "team_id", "identity_id", "role_id" FROM "team_identity_roles";
