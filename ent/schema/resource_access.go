package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ResourceAccess holds the schema definition for the ResourceAccess entity.
type ResourceAccess struct {
	ent.Schema
}

// Fields of the ResourceAccess.
func (ResourceAccess) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("tenant_id").Immutable(),
		field.String("resource_type").Immutable().Nillable().Optional(),
		field.Int64("resource_id").Immutable().Nillable().Optional(),
		field.String("identity_id").Immutable().Default(""),
		field.Int64("role_id").Immutable(),
	}
}

// Edges of the ResourceAccess.
func (ResourceAccess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Immutable().
			Required().
			Unique().
			Field("role_id"),
		edge.To("type", ResourceType.Type).
			Immutable().
			Unique().
			Field("resource_type"),
	}
}

// Indexes of the ResourceAccess.
func (ResourceAccess) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "role_id", "identity_id", "resource_type", "resource_id").
			Unique().
			Annotations(entsql.IndexWhere("resource_id IS NOT NULL")),
		index.Fields("tenant_id", "role_id", "identity_id").
			Unique().
			Annotations(entsql.IndexWhere("resource_id IS NULL")),
	}
}
