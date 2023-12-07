package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.Int64("tenant_id").Immutable().Nillable().Default(0),
		field.Int64("role_id").Immutable(),
		field.String("permission_id").Immutable(),
		field.Bool("deny").Default(false),
		field.JSON("fields", []string{}),
	}
}

// Edges of the RolePermission.
func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("permissions").
			Required().
			Unique().
			Field("role_id").Immutable(),
		edge.From("permission", Permission.Type).
			Ref("roles").
			Required().
			Unique().
			Field("permission_id").Immutable(),
	}
}

// Indexes of the RolePermission.
func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").Unique(),
	}
}
