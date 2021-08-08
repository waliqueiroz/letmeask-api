package entities

type Author struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Avatar string `json:"avatar" bson:"avatar"`
}
