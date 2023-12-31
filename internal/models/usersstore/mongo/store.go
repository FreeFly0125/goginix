package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

const (
	DBUsers          = "users"
	CollectionUsers  = "users"
	CollectionBadges = "badges"
)

type UsersStore struct {
	client *mongo.Client
}

// NewUsersStore makes connection to mongo server by provided url
// and return an instance of the client
func NewUsersStore(ctx context.Context, url string) (*UsersStore, error) {
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
	return &UsersStore{
		client: client,
	}, nil
}

func (s *UsersStore) getCollection(collection string) *mongo.Collection {
	return s.client.Database(DBUsers).Collection(collection)
}
