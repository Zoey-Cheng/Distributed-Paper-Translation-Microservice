package paper

import "time"

type Paper struct {
	ID             string    `bson:"ID"`
	FileHash       string    `bson:"FileHash"`
	CreateAt       time.Time `bson:"CreateAt"`
	Status         int32     `bson:"Status"`
	EmailTo        string    `bson:"EmailTo"`
	ResultText     string    `bson:"ResultText"`
	TargetLanguage string    `bson:"TargetLanguage"`
}
