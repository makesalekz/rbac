-- Create project permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id")
VALUES ('project.comment', 'Project Comment', 'common');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('project.task.create', 'project.task', 'Task Create', '', 'common', '[]'),
       ('project.task.read', 'project.task', 'Task Read', '', 'common', '[]'),
       ('project.task.update', 'project.task', 'Task Update', '', 'common', '[]'),
       ('project.task.delete', 'project.task', 'Task Delete', '', 'common', '[]'),
       ('project.comment.create', 'project.comment', 'Comment Create', '', 'common', '[]'),
       ('project.comment.read', 'project.comment', 'Comment Read', '', 'common', '[]'),
-- By default people can delete and update only there own comments,
-- except tenant owner, who can edit someone else's comments and delete them
-- and project manager, who can delete comments
       ('project.comment.update', 'project.project', 'Comment Update', '', 'common', '[]'),
       ('project.comment.delete', 'project.project', 'Comment Delete', '', 'common', '[]');

-- Create permissions for basic roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES
-- Tenant admin role
(1, 'project.task.create', FALSE, '[]'),
(1, 'project.task.read', FALSE, '[]'),
(1, 'project.task.update', FALSE, '[]'),
(1, 'project.task.delete', FALSE, '[]'),
(1, 'project.comment.create', FALSE, '[]'),
(1, 'project.comment.read', FALSE, '[]'),
(1, 'project.comment.update', FALSE, '[]'),
(1, 'project.comment.delete', FALSE, '[]'),
-- Project manager role
(3, 'project.task.read', FALSE, '[]'),
(3, 'project.task.update', FALSE, '[]'),
(3, 'project.task.delete', FALSE, '[]'),
(3, 'project.comment.create', FALSE, '[]'),
(3, 'project.comment.read', FALSE, '[]'),
(3, 'project.comment.delete', FALSE, '[]'),
-- Project participant limitations
(4, 'project.task.read', FALSE, '[]'),
(4, 'project.task.update', FALSE, '[]'),
(4, 'project.task.delete', FALSE, '[]'),
(4, 'project.comment.create', FALSE, '[]'),
(4, 'project.comment.read', FALSE, '[]'),
-- Project reporter limitations
(5, 'project.task.read', FALSE, '[]'),
(5, 'project.comment.create', FALSE, '[]'),
(5, 'project.comment.read', FALSE, '[]');