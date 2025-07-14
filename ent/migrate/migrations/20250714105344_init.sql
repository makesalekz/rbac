INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES (3, 'knowledge.document.create', FALSE, '[]'),
       (3, 'knowledge.document.read', FALSE, '[]'),
       (3, 'knowledge.document.update', FALSE, '[]'),
       (3, 'knowledge.document.delete', FALSE, '[]'),
       (3, 'knowledge.document.moderate', FALSE, '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES (4, 'knowledge.document.create', FALSE, '[]'),
       (4, 'knowledge.document.read', FALSE, '[]'),
       (4, 'knowledge.document.update', FALSE, '[]'),
       (4, 'knowledge.document.moderate', FALSE, '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields")
VALUES (5, 'knowledge.document.read', FALSE, '[]');