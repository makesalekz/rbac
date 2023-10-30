package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UsersRoles holds the schema definition for the UsersRoles entity.
type UsersRoles struct {
	ent.Schema
}

// Fields of the UsersRoles.
func (UsersRoles) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Int64("role_id"),
		field.Int64("team_id"),
	}
}

// Edges of the UsersRoles.
func (UsersRoles) Edges() []ent.Edge {
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

// Indexes of the UsersRoles
func (UsersRoles) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "role_id").Unique(),
	}
}
