package models

type Author struct {
	ID     string `bson:"_id,omitempty"`
	Name   string `bson:"name"`
	Avatar string `bson:"avatar"`
}
