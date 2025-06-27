package dbschema

import (
	"github.com/intezya/auth_service/internal/domain/access_level"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Account struct {
	ent.Schema
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.String("username").NotEmpty().Unique(),
		//field.String("email").Nillable().Optional().Unique(),
		field.String("password").NotEmpty().Sensitive(),
		field.String("hardware_id").Nillable().Optional().Unique().Sensitive(),

		field.String("access_level").
			GoType(access_level.AccessLevel(0)).
			DefaultFunc(
				func() access_level.AccessLevel {
					return access_level.User
				},
			),

		//field.String("avatar_url").Optional().Nillable(),
		//
		//field.Time("last_login_at").Default(time.Now),

		field.Time("created_at").Default(time.Now).Immutable(),

		field.Time("banned_until").Optional().Nillable(),
		field.String("ban_reason").Optional().Nillable(),
	}
}
