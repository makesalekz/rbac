-- Create resource permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id") VALUES
  ('resources.catalog', 'Resources Catalog', 'common'),
  ('resources.plan', 'Resources Plan', 'common');

-- Create resource permissions
INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('resources.catalog.create', 'resources.catalog', 'Resources Catalog Create', '', 'common', '[]'),
       ('resources.catalog.update', 'resources.catalog', 'Resources Catalog Update', '', 'common', '[]'),
       ('resources.catalog.read', 'resources.catalog', 'Resources Catalog Read', '', 'common', '[]'),
       ('resources.plan.create', 'resources.plan', 'Resources Plan Create', '', 'common', '[]'),
       ('resources.plan.update', 'resources.plan', 'Resources Plan Update', '', 'common', '[]'),
       ('resources.plan.read', 'resources.plan', 'Resources Plan Read', '', 'common', '[]'),
       ('resources.plan.delete', 'resources.plan', 'Resources Plan Delete', '', 'common', '[]');

-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES
-- Tenant admin role
(1, 'resources.catalog.create', FALSE, '[]'),
(1, 'resources.catalog.update', FALSE, '[]'),
(1, 'resources.catalog.read', FALSE, '[]'),
(1, 'resources.plan.create', FALSE, '[]'),
(1, 'resources.plan.update', FALSE, '[]'),
(1, 'resources.plan.read', FALSE, '[]'),
(1, 'resources.plan.delete', FALSE, '[]'),
-- Basic role
(2, 'resources.catalog.read', FALSE, '[]'),
(2, 'resources.plan.read', FALSE, '[]'),
-- Project manager role
(3, 'resources.catalog.create', FALSE, '[]'),
(3, 'resources.catalog.update', FALSE, '[]'),
(3, 'resources.catalog.read', FALSE, '[]'),
(3, 'resources.plan.create', FALSE, '[]'),
(3, 'resources.plan.update', FALSE, '[]'),
(3, 'resources.plan.read', FALSE, '[]'),
(3, 'resources.plan.delete', FALSE, '[]'),
-- Project viewer
(5, 'resources.catalog.read', FALSE, '[]'),
(5, 'resources.plan.read', FALSE, '[]');