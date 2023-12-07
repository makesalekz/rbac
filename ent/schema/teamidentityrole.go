package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TeamIdentityRole holds the schema definition for the TeamIdentityRole entity.
type TeamIdentityRole struct {
	ent.Schema
}

// Fields of the TeamIdentityRole.
func (TeamIdentityRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("tenant_id").Immutable(),
		field.Int64("team_id").Immutable().Nillable().Optional(),
		field.String("identity_id").Immutable().Default(""),
		field.Int64("role_id").Immutable(),
	}
}

// Edges of the TeamIdentityRole.
func (TeamIdentityRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Immutable().
			Required().
			Unique().
			Field("role_id"),
		edge.To("team", Team.Type).
			Immutable().
			Unique().
			Field("team_id"),
	}
}

// Indexes of the TeamIdentityRole.
func (TeamIdentityRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "role_id", "identity_id", "team_id").
			Unique().
			Annotations(entsql.IndexWhere("team_id IS NOT NULL")),
		index.Fields("tenant_id", "role_id", "identity_id").
			Unique().
			Annotations(entsql.IndexWhere("team_id IS NULL")),
	}
}
