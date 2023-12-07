package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("name").MaxLen(32).NotEmpty(),
		field.String("description").Optional().Default(""),
		field.String("app_id").MaxLen(10).Immutable(),
		field.JSON("fields", []string{}).Optional().Default([]string{}),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", RolePermission.Type),
	}
}
