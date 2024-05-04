-- Create project permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id") VALUES
  ('project.project', 'Project Manipulation', 'common'),
  ('project.team', 'Project Team', 'common'),
  ('project.attachments', 'Project Attachments', 'common');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields") VALUES
  ('project.project.create', 'project.project', 'Project Create', '','common', '[]'),
  ('project.project.read', 'project.project', 'Project Read', '','common', '[]'),
  ('project.project.update', 'project.project', 'Project Update', '','common', '[]'),
  ('project.project.delete', 'project.project', 'Project Delete', '','common', '[]'),
  ('project.project.manager', 'project.project', 'Project Manager Manipulation', '','common', '[]'),
  ('project.team.assign', 'project.team', 'Project Team Assign', '','common', '[]'),
  ('project.attachments.create','project.attachments','Project Attachments Create', '','common', '[]'),
  ('project.attachments.delete','project.attachments','Project Attachments Delete', '','common', '[]');

-- Create project roles
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at") VALUES
  (3, 'Project Owner', 'Has all permissions of a project owner', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (4, 'Project Manager', 'Has all permissions of a project manager', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (5, 'Project Participant', 'Has all limitations of project participant', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (6, 'Project Viewer', 'Has all limitations of project viewer', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

  (7, 'Direct Project Owner', 'Has all permissions of a direct project owner', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (8, 'Direct Project Manager', 'Has all permissions of a direct project manager', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (9, 'Direct Project Participant', 'Has all limitations of direct project participant', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (10, 'Direct Project Viewer', 'Has all limitations of direct project viewer', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields") VALUES
-- Tenant admin role
  (1, 'project.project.create', FALSE, '[]'),
  (1, 'project.project.read', FALSE, '[]'),
  (1, 'project.project.update', FALSE, '[]'),
  (1, 'project.project.delete', FALSE, '[]'),
  (1, 'project.project.manager', FALSE, '[]'),
  (1, 'project.team.assign', FALSE, '[]'),
  (1, 'project.attachments.create',FALSE,'[]'),
  (1, 'project.attachments.delete',FALSE,'[]'),
-- Basic plan limitations
  (2, 'project.project.create', FALSE, '[]'),
  (2, 'project.project.read', FALSE, '[]'),
  (2, 'project.project.update', TRUE, '[]'),
  (2, 'project.project.delete', TRUE, '[]'),
  (2, 'project.project.manager', TRUE, '[]'),
  (2, 'project.team.assign', TRUE, '[]'),
  (2, 'project.attachments.create',TRUE,'[]'),
  (2, 'project.attachments.delete',TRUE,'[]'),
-- Project owner role
  (3, 'admin.team.read', FALSE, '[]'),
  (3, 'project.project.create', FALSE, '[]'),
  (3, 'project.project.read', FALSE, '[]'),
  (3, 'project.project.update', FALSE, '[]'),
  (3, 'project.project.delete', FALSE, '[]'),
  (3, 'project.project.manager', FALSE, '[]'),
  (3, 'project.team.assign', FALSE, '[]'),
  (3, 'project.attachments.create',FALSE,'[]'),
  (3, 'project.attachments.delete',FALSE,'[]'),
-- Project manager role
  (4, 'admin.team.read', FALSE, '[]'),
  (4, 'project.project.create', FALSE, '[]'),
  (4, 'project.project.read', FALSE, '[]'),
  (4, 'project.project.update', FALSE, '[]'),
  (4, 'project.project.delete', FALSE, '[]'),
  (4, 'project.project.manager', TRUE, '[]'),
  (4, 'project.team.assign', FALSE, '[]'),
  (4, 'project.attachments.create',FALSE,'[]'),
  (4, 'project.attachments.delete',FALSE,'[]'),
-- Project participant limitations
  (5, 'admin.team.read', TRUE, '[]'),
  (5, 'project.project.create', FALSE, '[]'),
  (5, 'project.project.read', FALSE, '[]'),
  (5, 'project.project.update', TRUE, '[]'),
  (5, 'project.project.delete', TRUE, '[]'),
  (5, 'project.project.manager', TRUE, '[]'),
  (5, 'project.team.assign', TRUE, '[]'),
  (5, 'project.attachments.create',FALSE,'[]'),
  (5, 'project.attachments.delete',FALSE,'[]'),
-- Project reporter limitations
  (6, 'admin.team.read', TRUE, '[]'),
  (6, 'project.project.create', TRUE, '[]'),
  (6, 'project.project.read', FALSE, '[]'),
  (6, 'project.project.update', TRUE, '[]'),
  (6, 'project.project.delete', TRUE, '[]'),
  (6, 'project.project.manager', TRUE, '[]'),
  (6, 'project.team.assign', TRUE, '[]'),
  (6, 'project.attachments.create',TRUE,'[]'),
  (6, 'project.attachments.delete',TRUE,'[]'),

-- Project direct owner role
  (7, 'admin.team.read', FALSE, '[]'),
  (7, 'project.project.create', FALSE, '[]'),
  (7, 'project.project.read', TRUE, '[]'),
  (7, 'project.project.update', FALSE, '[]'),
  (7, 'project.project.delete', FALSE, '[]'),
  (7, 'project.project.manager', FALSE, '[]'),
  (7, 'project.team.assign', FALSE, '[]'),
  (7, 'project.attachments.create',FALSE,'[]'),
  (7, 'project.attachments.delete',FALSE,'[]'),
-- Project direct manager role
  (8, 'admin.team.read', FALSE, '[]'),
  (8, 'project.project.create', FALSE, '[]'),
  (8, 'project.project.read', TRUE, '[]'),
  (8, 'project.project.update', FALSE, '[]'),
  (8, 'project.project.delete', FALSE, '[]'),
  (8, 'project.project.manager', TRUE, '[]'),
  (8, 'project.team.assign', FALSE, '[]'),
  (8, 'project.attachments.create',FALSE,'[]'),
  (8, 'project.attachments.delete',FALSE,'[]'),
-- Project direct participant limitations
  (9, 'admin.team.read', TRUE, '[]'),
  (9, 'project.project.create', FALSE, '[]'),
  (9, 'project.project.read', TRUE, '[]'),
  (9, 'project.project.update', TRUE, '[]'),
  (9, 'project.project.delete', TRUE, '[]'),
  (9, 'project.project.manager', TRUE, '[]'),
  (9, 'project.team.assign', TRUE, '[]'),
  (9, 'project.attachments.create',FALSE,'[]'),
  (9, 'project.attachments.delete',FALSE,'[]'),
-- Project direct reporter limitations
  (10, 'admin.team.read', TRUE, '[]'),
  (10, 'project.project.create', TRUE, '[]'),
  (10, 'project.project.read', TRUE, '[]'),
  (10, 'project.project.list', TRUE, '[]'),
  (10, 'project.project.update', TRUE, '[]'),
  (10, 'project.project.delete', TRUE, '[]'),
  (10, 'project.project.manager', TRUE, '[]'),
  (10, 'project.team.assign', TRUE, '[]'),
  (10, 'project.attachments.create',TRUE,'[]'),
  (10, 'project.attachments.delete',TRUE,'[]');
