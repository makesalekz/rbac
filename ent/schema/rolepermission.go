package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RolePermission holds the schema definition for the RolePermission entity.
type RolePermission struct {
	ent.Schema
}

// Fields of the RolePermission.
func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("role_id").Immutable(),
		field.Int64("permission_id").Immutable(),
		field.JSON("fields", []string{}),
	}
}

// Edges of the RolePermission.
func (RolePermission) Edges() []ent.Edge {
	return nil
}

// Indexes of the RolePermission.
func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").Unique(),
	}
}
