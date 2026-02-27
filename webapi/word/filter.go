package word

import "go.mongodb.org/mongo-driver/v2/bson"

type WordFilter struct {
	Difficulty string
}

func (f WordFilter) ToBSON() bson.D {
	filter := bson.D{}

	if f.Difficulty != "" {
		filter = append(filter, bson.E{Key: "difficulty", Value: f.Difficulty})
	}

	return filter
}
