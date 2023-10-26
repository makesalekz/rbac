package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Immutable(),
		field.String("key").MaxLen(32).NotEmpty().Unique(),
		field.String("name"),
		field.Int64("organization_id"),
		field.Int64("app_id"),
		field.JSON("fields", []string{}),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type),
	}
}
