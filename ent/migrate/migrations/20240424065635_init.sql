-- Create project permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id") VALUES
  ('project.project', 'Project Manipulation', 'common'),
  ('project.team', 'Project Team', 'common'),
  ('project.attachment', 'Project Attachment', 'common');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields") VALUES
  ('project.project.create', 'project.project', 'Project Create', '','common', '[]'),
  ('project.project.read', 'project.project', 'Project Read', '','common', '[]'),
  ('project.project.update', 'project.project', 'Project Update', '','common', '[]'),
  ('project.project.delete', 'project.project', 'Project Delete', '','common', '[]'),
  ('project.project.manager', 'project.project', 'Project Manager Manipulation', '','common', '[]'),
  ('project.team.assign', 'project.team', 'Project Team Assign', '','common', '[]'),
  ('project.attachment.create','project.attachment','Project Attachment Create', '','common', '[]'),
  ('project.attachment.delete','project.attachment','Project Attachment Delete', '','common', '[]');

-- Create project roles
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at") VALUES
  (3, 'Project Manager', 'Has all permissions of a project manager', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (4, 'Project Participant', 'Has all limitations of project participant', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (5, 'Project Viewer', 'Has all limitations of project viewer', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

  (6, 'Direct Participant', 'Has all permissions of a Direct Participant', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields") VALUES
-- Tenant admin role
  (1, 'project.project.create', FALSE, '[]'),
  (1, 'project.project.read', FALSE, '[]'),
  (1, 'project.project.update', FALSE, '[]'),
  (1, 'project.project.delete', FALSE, '[]'),
  (1, 'project.project.manager', FALSE, '[]'),
  (1, 'project.team.assign', FALSE, '[]'),
  (1, 'project.attachment.create',FALSE,'[]'),
  (1, 'project.attachment.delete',FALSE,'[]'),
-- Project manager role
  (3, 'admin.team.read', FALSE, '[]'),
  (3, 'project.project.create', FALSE, '[]'),
  (3, 'project.project.read', FALSE, '[]'),
  (3, 'project.project.update', FALSE, '[]'),
  (3, 'project.project.delete', FALSE, '[]'),
  (3, 'project.team.assign', FALSE, '[]'),
  (3, 'project.attachment.create',FALSE,'[]'),
  (3, 'project.attachment.delete',FALSE,'[]'),
-- Project participant limitations
  (4, 'project.project.create', FALSE, '[]'),
  (4, 'project.project.read', FALSE, '[]'),
  (4, 'project.attachment.create',FALSE,'[]'),
  (4, 'project.attachment.delete',FALSE,'[]'),
-- Project reporter limitations
  (5, 'project.project.create', FALSE, '[]'),
  (5, 'project.project.read', FALSE, '[]'),

-- Direct particiant limitation
  (6, 'project.project.create', FALSE, '[]');