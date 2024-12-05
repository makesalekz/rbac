-- Create resource permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id") VALUES
  ('project.catalog', 'Resources Catalog', 'common'),
  ('project.plan', 'Resources Plan', 'common');

-- Create resource permissions
INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('project.catalog.create', 'project.catalog', 'Resources Catalog Create', '', 'common', '[]'),
       ('project.catalog.update', 'project.catalog', 'Resources Catalog Update', '', 'common', '[]'),
       ('project.catalog.read', 'project.catalog', 'Resources Catalog Read', '', 'common', '[]'),
       ('project.plan.create', 'project.plan', 'Resources Plan Create', '', 'common', '[]'),
       ('project.plan.update', 'project.plan', 'Resources Plan Update', '', 'common', '[]'),
       ('project.plan.read', 'project.plan', 'Resources Plan Read', '', 'common', '[]'),
       ('project.plan.delete', 'project.plan', 'Resources Plan Delete', '', 'common', '[]');

-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES
-- Tenant admin role
(1, 'project.catalog.create', FALSE, '[]'),
(1, 'project.catalog.update', FALSE, '[]'),
(1, 'project.catalog.read', FALSE, '[]'),
(1, 'project.plan.create', FALSE, '[]'),
(1, 'project.plan.update', FALSE, '[]'),
(1, 'project.plan.read', FALSE, '[]'),
(1, 'project.plan.delete', FALSE, '[]'),
-- Basic role
(2, 'project.catalog.read', FALSE, '[]'),
(2, 'project.plan.read', FALSE, '[]'),
-- Project manager role
(3, 'project.catalog.create', FALSE, '[]'),
(3, 'project.catalog.update', FALSE, '[]'),
(3, 'project.catalog.read', FALSE, '[]'),
(3, 'project.plan.create', FALSE, '[]'),
(3, 'project.plan.update', FALSE, '[]'),
(3, 'project.plan.read', FALSE, '[]'),
(3, 'project.plan.delete', FALSE, '[]'),
-- Project viewer
(5, 'project.catalog.read', FALSE, '[]'),
(5, 'project.plan.read', FALSE, '[]');