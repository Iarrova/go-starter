package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/Iarrova/go-starter/ent/schema/mixins"
	"github.com/google/uuid"
)

type Post struct {
	ent.Schema
}

func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID{},
		mixins.Timestamps{},
	}
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}).
			Immutable(),

		field.String("title").
			NotEmpty().
			MaxLen(200),

		field.Text("body").
			NotEmpty(),
	}
}

func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("posts").
			Field("user_id").
			Immutable().
			Unique().
			Required(),
	}
}

func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "created_at"),
	}
}

func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Table("posts"),
	}
}
