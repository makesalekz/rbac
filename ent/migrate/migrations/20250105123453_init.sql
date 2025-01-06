INSERT INTO "permission_groups" ("id", "name", "app_id")
VALUES ('basqaru.limits', 'Basqaru Limits', 'pms');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('basqaru.limits.media', 'basqaru.limits', 'Limits for Media Size', '', 'pms', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (1, 'basqaru.limits.media', FALSE, '[]', 0);

-- Create basic roles
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at")
VALUES (9, 'Basqaru Basic User', 'Has all limits of a Basqaru Basic User', TRUE, CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (9, 'basqaru.limits.media', TRUE, '[]', 20 * 1000 * 1000);

INSERT INTO resource_accesses(tenant_id, resource_id, identity_id, role_id, resource_type)
SELECT tenant_id, null, '', 9, null
FROM resource_accesses
GROUP BY tenant_id
HAVING SUM(CASE WHEN role_id = 9 THEN 1 ELSE 0 END) = 0;

