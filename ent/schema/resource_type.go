package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ResourceType holds the schema definition for the ResourceType entity.
type ResourceType struct {
	ent.Schema
}

// Fields of the Resource.
func (ResourceType) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("description").Optional().Default(""),
	}
}

// Edges of the Resource.
func (ResourceType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type),
	}
}
