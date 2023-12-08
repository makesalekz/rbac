package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PermissionGroup holds the schema definition for the PermissionGroup entity.
type PermissionGroup struct {
	ent.Schema
}

// Fields of the PermissionGroup.
func (PermissionGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("app_id").MaxLen(10).Immutable(),
		field.String("name").MaxLen(16).NotEmpty(),
	}
}

// Edges of the PermissionGroup.
func (PermissionGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("permissions", Permission.Type),
	}
}
