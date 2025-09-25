package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// User — Ent-схема пользователя (модель находится в internal/models).
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{}, // created_at, updated_at
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default(""),
		field.String("email").NotEmpty().Unique(),
	}
}
