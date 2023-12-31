package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

const (
	DBStories           = "stories"
	CollectionStories   = "stories"
	CollectionComments  = "comments"
	CollectionReactions = "reactions"
)

type StoriesStore struct {
	client *mongo.Client
}

// NewStoriesStore makes connection to mongo server by provided url
// and return an instance of the client
func NewStoriesStore(ctx context.Context, url string) (*StoriesStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url).SetMonitor(otelmongo.NewMonitor()))
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &StoriesStore{
		client: client,
	}, nil
}

func (s *StoriesStore) getCollection(collection string) *mongo.Collection {
	return s.client.Database(DBStories).Collection(collection)
}
