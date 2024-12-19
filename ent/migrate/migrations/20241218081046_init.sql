INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at")
VALUES (8, 'Qalai Trial User', 'Has all limits of a Qalai Trial User', TRUE, CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP);

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('qalai.limits.recap', 'qalai.limits', 'Limits for recap duration', '', 'calendaria', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (1, 'qalai.limits.recap', FALSE, '[]', 0);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (7, 'qalai.limits.recap', TRUE, '[]', 0);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (8, 'qalai.limits.recap', TRUE, '[]', 30 * 60); -- 30 minutes

INSERT INTO resource_accesses(tenant_id, resource_id, identity_id, role_id, resource_type)
SELECT tenant_id, null, '', 8, null
FROM resource_accesses
GROUP BY tenant_id
HAVING SUM(CASE WHEN role_id = 8 THEN 1 ELSE 0 END) = 0;

