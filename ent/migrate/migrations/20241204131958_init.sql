-- Modify "role_permissions" table
ALTER TABLE "role_permissions"
    ADD COLUMN "value" bigint NULL DEFAULT 0;

-- Allow project manager to assign roles
INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES
-- Project manager role
(3, 'admin.role.assign', FALSE, '[]');

-- Create basic roles
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at")
VALUES (7, 'Qalai Basic User', 'Has all limits of a Qalai Basic User', TRUE, CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP);

-- Create qalai permission groups
INSERT INTO "permission_groups" ("id", "name", "app_id")
VALUES ('qalai.limits', 'Qalai Limits', 'calendaria');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('qalai.limits.media', 'qalai.limits', 'Limits for Media Size', '', 'calendaria', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (1, 'qalai.limits.media', FALSE, '[]', 0);

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields", "value")
VALUES (7, 'qalai.limits.media', TRUE, '[]', 100 * 1000 * 1000);