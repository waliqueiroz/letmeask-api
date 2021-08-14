package entities

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
