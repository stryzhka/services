package models

type Word struct {
	Id              string `bson:"_id,omitempty"    json:"id"`
	EngText         string `bson:"eng_text"         json:"eng_text"`
	NativeText      string `bson:"native_text"      json:"native_text"`
	Transcription   string `bson:"transcription"    json:"transcription"`
	Difficulty      string `bson:"difficulty"       json:"difficulty"`
	CategoryId      string `bson:"category_id"      json:"category_id"`
	ConfirmedUserId string `bson:"confirmed_user_id" json:"confirmed_user_id"`
	ConfirmedStatus string `bson:"confirmed_status" json:"confirmed_status"`
	ConfirmedAt     string `bson:"confirmed_at" json:"confirmed_at"`
}
