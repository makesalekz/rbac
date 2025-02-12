--  Create roles for superadmin
INSERT INTO "roles" ("id", "name", "description", "is_system", "created_at", "updated_at")
VALUES (10, 'Administrator', 'System configuration, management, monitoring', TRUE, CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP),
       (11, 'Manager', 'Manages organizers, invites, roles, finances, reports', TRUE, CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP);

-- Add permissions for organizations
INSERT INTO "permission_groups" ("id", "name", "app_id")
VALUES ('admin.organization', 'Organization', 'admin');

INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields")
VALUES ('admin.organization.create', 'admin.organization', 'Create organization', '', 'admin', '[]'),
       ('admin.organization.update', 'admin.organization', 'Update organization', '', 'admin', '[]'),
       ('admin.organization.delete', 'admin.organization', 'Delete organization', '', 'admin', '[]'),
       ('admin.organization.read', 'admin.organization', 'Get organization', '', 'admin', '[]'),
       ('admin.legal.create', 'admin.organization', 'Create organization legal', '', 'admin', '[]'),
       ('admin.legal.update', 'admin.organization', 'Update organization legal', '', 'admin', '[]'),
       ('admin.legal.delete', 'admin.organization', 'Delete organization legal', '', 'admin', '[]'),
       ('admin.legal.read', 'admin.organization', 'Get organization legal', '', 'admin', '[]'),
       ('admin.legal_contact.create', 'admin.organization', 'Create legal contact', '', 'admin', '[]'),
       ('admin.legal_contact.update', 'admin.organization', 'Update legal contact', '', 'admin', '[]'),
       ('admin.legal_contact.delete', 'admin.organization', 'Delete legal contact', '', 'admin', '[]'),
       ('admin.legal_contact.read', 'admin.organization', 'Delete legal contact', '', 'admin', '[]');

--  Give permissions to the admin role
insert into role_permissions (role_id, permission_id, deny, fields, value)
        (1, 'admin.organization.create', FALSE, '[]', 0),
    (1, 'admin.organization.update', FALSE, '[]', 0),
    (1, 'admin.organization.delete', FALSE, '[]', 0),
    (1, 'admin.organization.get', FALSE, '[]', 0),
    (1, 'admin.legal.create', FALSE, '[]', 0),
    (1, 'admin.legal.update', FALSE, '[]', 0),
    (1, 'admin.legal.delete', FALSE, '[]', 0),
    (1, 'admin.legal.get', FALSE, '[]', 0),
    (1, 'admin.legal_contact.create', FALSE, '[]', 0),
    (1, 'admin.legal_contact.update', FALSE, '[]', 0),
    (1, 'admin.legal_contact.delete', FALSE, '[]', 0),
    (1, 'admin.legal_contact.get', FALSE, '[]', 0);

-- Copy permissions from the admin role to the superadmin role
insert into role_permissions (role_id, permission_id, deny, fields, value)
        (SELECT 10, id, FALSE, '[]', 0 FROM permissions WHERE app_id = 'admin'),
     (SELECT 11, id, FALSE, '[]', 0 FROM permissions WHERE app_id = 'admin'
          AND id NOT LIKE '%.invite.%'
          AND id <> 'admin.tenant.delete'
          AND id NOT LIKE '%.role_system.%'
          AND id NOT LIKE '%role%';

     );

