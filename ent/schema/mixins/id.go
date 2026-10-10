package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type ID struct {
	mixin.Schema
}

func (ID) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(func() uuid.UUID {
				return uuid.Must(uuid.NewV7())
			}),
	}
}
