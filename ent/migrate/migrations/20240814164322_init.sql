INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('project.attachment.read', 'project.attachment', 'Project Attachment Read', '', 'common', '[]');

-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES
-- Tenant admin role
(1, 'project.attachment.read', FALSE, '[]'),
-- Project manager role
(3, 'project.attachment.read', FALSE, '[]'),
-- Project participant limitations
(4, 'project.attachment.read', FALSE, '[]'),
-- Project reporter limitations
(5, 'project.attachment.read', FALSE, '[]');