package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberRole holds the schema definition for the MemberRole entity.
type MemberRole struct {
	ent.Schema
}

// Fields of the MemberRole.
func (MemberRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Immutable().Unique(),
		field.Int64("member_id").Immutable(),
		field.Int64("role_id"),
		field.Int64("team_id"),
	}
}

// Edges of the MemberRole.
func (MemberRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Required().
			Unique().
			Field("role_id"),
		edge.To("teams", Role.Type).
			Required().
			Unique().
			Field("team_id"),
	}
}

func (MemberRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("member_id", "role_id", "team_id").Unique(),
	}
}
