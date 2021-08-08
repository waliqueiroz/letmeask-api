package entities

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	safeUser := struct {
		Password string `json:"password,omitempty"`
		Alias
	}{
		Alias: Alias(u),
	}

	return json.Marshal(safeUser)
}
