package schema

import (
	"time"

	"gitlab.calendaria.team/services/rbac/ent/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/jackc/pgtype"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("tenant_id"),
		field.Int64("parent_id").Nillable().Optional(),
		field.Other("parents_ids", &pgtype.Int8Array{}).
			SchemaType(map[string]string{
				dialect.Postgres: "bigint[]",
			}).Optional(),
		field.String("name"),
		field.String("description").Default("").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Team.Type).
			From("parent").
			Unique().
			Field("parent_id"),
	}
}

func (Team) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.SoftDeleteMixin{},
	}
}
