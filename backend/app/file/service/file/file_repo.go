package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	Query(hash string) (*File, error)
	Create(file *File) error
	Update(hash string, set map[string]any) error
	Delete(hash string) error
}

type MongoFileRepository struct {
	C *mongo.Collection
}

func NewMongoFileRepository(db *mongo.Database) FileRepository {
	return &MongoFileRepository{C: db.Collection("files")}
}

func (t *MongoFileRepository) Query(hash string) (f *File, err error) {
	return f, t.C.FindOne(context.TODO(), bson.M{"Hash": hash}).Decode(&f)
}

func (t *MongoFileRepository) Create(file *File) error {
	_, err := t.C.InsertOne(context.TODO(), file)
	return err
}

func (t *MongoFileRepository) Update(hash string, set map[string]any) error {
	_, err := t.C.UpdateOne(context.TODO(), bson.M{"Hash": hash}, bson.M{"$set": set})
	return err
}

func (t *MongoFileRepository) Delete(hash string) error {
	_, err := t.C.DeleteOne(context.TODO(), bson.M{"Hash": hash})
	return err
}
