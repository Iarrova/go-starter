package mixins

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Timestamps struct {
	mixin.Schema
}

func (Timestamps) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(entsql.Default("CURRENT_TIMESTAMP")),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entsql.Default("CURRENT_TIMESTAMP")),
	}
}
