package models

// string поменяется на uuid впоследствии. memdb не умеет в uuid
type Word struct {
	Id            string `bson:"_id,omitempty"    json:"id"`
	EngText       string `bson:"eng_text"         json:"eng_text"`
	NativeText    string `bson:"native_text"      json:"native_text"`
	Transcription string `bson:"transcription"    json:"transcription"`
	Difficulty    string `bson:"difficulty"       json:"difficulty"`
	CategoryId    string `bson:"category_id"      json:"category_id"`
}
