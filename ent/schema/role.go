package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.calendaria.team/services/rbac/ent/mixins"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Immutable(),
		field.String("name").MaxLen(32).NotEmpty(),
		field.String("description").Optional().Default(""),
		field.Int64("tenant_id").Nillable(),
		field.Bool("is_system").Default(false),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("permissions", RolePermission.Type),
	}
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.SoftDeleteMixin{},
	}
}
