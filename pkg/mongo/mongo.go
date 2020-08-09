package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

// Connection ...
type Connection struct {
	URL    string
	DBName string
}

// Connect ...
func (m Connection) Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(m.URL).SetWriteConcern(writeconcern.New(writeconcern.W(1), writeconcern.J(false)))
	client, err := mongo.Connect(ctx, clientOptions)

	return client, err
}
