package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("tenant_id"),
		field.String("name").NotEmpty(),
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return nil
}
