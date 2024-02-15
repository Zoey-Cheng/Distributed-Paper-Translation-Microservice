package file

type Chunk struct {
	ChunkIndex int64 `bson:"ChunkIndex"`
	ChunkOK    bool  `bson:"ChunkOK"`
}

type File struct {
	Hash         string  `bson:"Hash"`
	Status       int32   `bson:"Status"`
	ChunkNums    int64   `bson:"ChunkNums"`
	CurrentIndex int64   `bson:"CurrentIndex"`
	SegmentSize  int64   `bson:"SegmentSize"`
	Bucket       string  `bson:"Bucket"`
	FilePath     string  `bson:"FilePath"`
	Chunks       []Chunk `bson:"Chunks"`
}
