package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TeamIdentityRole holds the schema definition for the TeamIdentityRole entity.
type TeamIdentityRole struct {
	ent.Schema
}

// Fields of the TeamIdentityRole.
func (TeamIdentityRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("tenant_id").Immutable(),
		field.Int64("team_id").Immutable().Optional().Default(0),
		field.String("identity_id").Immutable().Default("0"),
		field.Int64("role_id").Immutable(),
	}
}

// Edges of the TeamIdentityRole.
func (TeamIdentityRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Required().
			Unique().
			Field("role_id").Immutable(),
		edge.To("team", Team.Type).
			Unique().
			Field("team_id").Immutable(),
	}
}
