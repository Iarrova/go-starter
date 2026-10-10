package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Iarrova/go-starter/ent/schema/mixins"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.ID{},
		mixins.Timestamps{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			NotEmpty().
			Unique(),

		field.String("password_hash").
			NotEmpty().
			Sensitive(),

		field.String("name").
			NotEmpty().
			MaxLen(200),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Table("users"),
	}
}
