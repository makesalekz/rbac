-- Fix for knowledge document permissions
UPDATE "permission_groups"
  SET "name" = 'Knowledge Document'
  WHERE "id" = 'knowledge.document';

UPDATE "permissions"
  SET "name" = 'Knowledge Document Create'
  WHERE "id" = 'knowledge.document.create';

UPDATE "permissions"
  SET "name" = 'Knowledge Document Read'
  WHERE "id" = 'knowledge.document.read';

UPDATE "permissions"
  SET "name" = 'Knowledge Document Update'
  WHERE "id" = 'knowledge.document.update';

UPDATE "permissions"
  SET "name" = 'Knowledge Document Delete'
  WHERE "id" = 'knowledge.document.delete';

UPDATE "permissions"
  SET "name" = 'Knowledge Document Moderate'
  WHERE "id" = 'knowledge.document.moderate';

-- Fix for project permissions
UPDATE "permission_groups" SET  "name" = 'Project' WHERE "id" = 'project.project';
DELETE FROM "role_permissions" WHERE "permission_id" = 'project.project.manager';
DELETE FROM "permissions" WHERE "id" = 'project.project.manager';
