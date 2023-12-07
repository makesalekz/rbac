-- Modify "roles" table
ALTER TABLE "roles" ALTER COLUMN "tenant_id" SET NOT NULL, ALTER COLUMN "tenant_id" SET DEFAULT 0;
-- Modify "team_identity_roles" table
ALTER TABLE "team_identity_roles" ALTER COLUMN "team_id" DROP DEFAULT;
-- Create index "teamidentityrole_tenant_id_role_id_identity_id" to table: "team_identity_roles"
CREATE UNIQUE INDEX "teamidentityrole_tenant_id_role_id_identity_id" ON "team_identity_roles" ("tenant_id", "role_id", "identity_id") WHERE (team_id IS NULL);
-- Create index "teamidentityrole_tenant_id_role_id_identity_id_team_id" to table: "team_identity_roles"
CREATE UNIQUE INDEX "teamidentityrole_tenant_id_role_id_identity_id_team_id" ON "team_identity_roles" ("tenant_id", "role_id", "identity_id", "team_id") WHERE (team_id IS NOT NULL);
