INSERT INTO "permission_groups" ("id", "name", "app_id") VALUES
  ('knowledge.document', 'Knowledge Documents', 'knowledge');
  
INSERT INTO "permissions" ("id", "group_id", "name", "description", "app_id", "fields") VALUES
  ('knowledge.document.create', 'knowledge.document', 'Document Create', '', 'knowledge', '[]'),
  ('knowledge.document.read', 'knowledge.document', 'Document Read', '', 'knowledge', '[]'),
  ('knowledge.document.update', 'knowledge.document', 'Document Update', '', 'knowledge', '[]'),
  ('knowledge.document.delete', 'knowledge.document', 'Document Delete', '', 'knowledge', '[]'),
  ('knowledge.document.moderate', 'knowledge.document', 'Document Moderate', '', 'knowledge', '[]');

INSERT INTO "role_permissions" ("role_id", "permission_id", "deny", "fields") VALUES
  (1, 'knowledge.document.create', FALSE, '[]'),
  (1, 'knowledge.document.read', FALSE, '[]'),
  (1, 'knowledge.document.update', FALSE, '[]'),
  (1, 'knowledge.document.delete', FALSE, '[]'),
  (1, 'knowledge.document.moderate', FALSE, '[]');