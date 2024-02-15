package ocr

type OCR struct {
	ID        string `bson:"ID"`
	Bucket    string `bson:"Bucket"`
	ObjectKey string `bson:"ObjectKey"`
	FileType  string `bson:"FileType"`
	OcredText string `bson:"OcredText"`
}
