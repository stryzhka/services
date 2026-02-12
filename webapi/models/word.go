package models

// string поменяется на uuid впоследствии. memdb не умеет в uuid
type Word struct {
	Id            string `json:"id"`
	EngText       string `json:"eng_text"`
	NativeText    string `json:"native_text"`
	Transcription string `json:"transcription"`
	Difficulty    string `json:"difficulty"`
	CategoryId    string `json:"category_id"`
}
