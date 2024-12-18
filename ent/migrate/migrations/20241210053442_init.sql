INSERT INTO resource_accesses(tenant_id, resource_id, identity_id, role_id, resource_type)
SELECT tenant_id, null, '', 7, null
FROM resource_accesses
GROUP BY tenant_id
HAVING SUM(CASE WHEN role_id = 7 THEN 1 ELSE 0 END) = 0;

