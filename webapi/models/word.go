package models

import "github.com/google/uuid"

type Word struct {
	Id            uuid.UUID `json:"id"`
	EngText       string    `json:"eng_text"`
	NativeText    string    `json:"native_text"`
	Transcription string    `json:"transcription"`
	Difficulty    string    `json:"difficulty"`
	CategoryId    uuid.UUID `json:"category_id"`
}
