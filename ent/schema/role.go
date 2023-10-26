package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Immutable(),
		field.String("key").MaxLen(16).NotEmpty().Unique(),
		field.String("name"),
		field.Int64("organization_id"),
		field.Int64("app_id"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("permissions", Permission.Type).
			Ref("roles"),
	}
}

// Indexes of the Role
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key"),
	}
}
