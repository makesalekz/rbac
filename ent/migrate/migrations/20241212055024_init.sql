-- Create qalai permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id")
VALUES ('qalai.features', 'Qalai Features', 'calendaria');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('qalai.features.event_pending_invite', 'qalai.features', 'Feature for Event Pending Invite', '', 'calendaria', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (1, 'qalai.features.event_pending_invite', FALSE, '[]', 0);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (7, 'qalai.features.event_pending_invite', TRUE, '[]', 0);