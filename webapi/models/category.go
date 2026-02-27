package models

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"` //потом поправлю
	Name string `bson:"name"           json:"name"`
}
