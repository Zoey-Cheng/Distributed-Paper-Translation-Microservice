package ocr

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OCRRepository interface {
	Create(ocr *OCR) error
	Get(bucket string, objectKey string, fileType string) (*OCR, error)
}

type MongoOCRRepository struct {
	C *mongo.Collection
}

func NewMongoOCRRepository(db *mongo.Database) *MongoOCRRepository {
	return &MongoOCRRepository{C: db.Collection("ocrs")}
}

func (t *MongoOCRRepository) Create(ocr *OCR) error {
	_, err := t.C.InsertOne(context.TODO(), ocr)
	return err
}

func (t *MongoOCRRepository) Get(bucket string, objectKey string, fileType string) (o *OCR, err error) {
	return o, t.C.FindOne(context.TODO(), bson.M{
		"Bucket":    bucket,
		"ObjectKey": objectKey,
		"FileType":  fileType,
	}).Decode(&o)
}
