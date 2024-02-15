package paper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaperRepository interface {
	Create(paper *Paper) error
	Get(id string) (*Paper, error)
	UpdateText(id string, text string) error
	SetStatus(id string, status int32) error
	Delete(id string) error
	GetPapers() ([]*Paper, error)
}

type MongoPaperRepository struct {
	C *mongo.Collection
}

func NewMongoPaperRepository(db *mongo.Database) *MongoPaperRepository {
	return &MongoPaperRepository{C: db.Collection("papers")}
}

func (t *MongoPaperRepository) Create(paper *Paper) error {
	_, err := t.C.InsertOne(context.Background(), paper)
	return err
}

func (t *MongoPaperRepository) Get(id string) (p *Paper, err error) {
	return p, t.C.FindOne(context.TODO(), bson.M{"ID": id}).Decode(&p)
}

func (t *MongoPaperRepository) UpdateText(id string, text string) error {
	_, err := t.C.UpdateOne(context.TODO(), bson.M{"ID": id}, bson.M{
		"$set": bson.M{
			"ResultText": text,
		},
	})
	return err
}

func (t *MongoPaperRepository) SetStatus(id string, status int32) error {
	_, err := t.C.UpdateOne(context.TODO(), bson.M{"ID": id}, bson.M{
		"$set": bson.M{
			"Status": status,
		},
	})
	return err
}

func (t *MongoPaperRepository) Delete(id string) error {
	_, err := t.C.DeleteOne(context.TODO(), bson.M{"ID": id})
	return err
}

func (t *MongoPaperRepository) GetPapers() (ps []*Paper, err error) {
	cur, err := t.C.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"CreateAt": -1}))
	if err != nil {
		return nil, err
	}
	return ps, cur.All(context.TODO(), &ps)
}
