INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('qalai.limits.aigenda_daily_request', 'qalai.limits', 'Limits for Media Size', '', 'calendaria', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (1, 'qalai.limits.aigenda_daily_request', FALSE, '[]', 0);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (7, 'qalai.limits.aigenda_daily_request', TRUE, '[]', 5);