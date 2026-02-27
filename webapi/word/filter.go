package word

import "go.mongodb.org/mongo-driver/v2/bson"

type WordFilter struct {
	Difficulty string
}

func (f WordFilter) ToBSON() bson.D {
	filter := bson.D{}

	if f.Difficulty != "" {
		filter = append(filter, bson.E{"difficulty", bson.D{{"$regex", f.Difficulty}, {"$options", "i"}}})
	}

	return filter
}
