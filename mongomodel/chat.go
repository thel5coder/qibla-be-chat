package mongomodel

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// ChatTypeText ...
	ChatTypeText = "text"
	// ChatTypeImage ...
	ChatTypeImage = "image"
	// ChatTypeVideo ...
	ChatTypeVideo = "video"
	// ChatTypeAudio ...
	ChatTypeAudio = "audio"
	// ChatTypeFile ...
	ChatTypeFile = "file"
)

// IChat ...
type IChat interface {
	FindAllByRoom(roomID, message, lastID string, limit int) ([]ChatEntity, error)
	FindLast(roomID string) (ChatEntity, error)
	Store(body *ChatEntity) (string, error)
	Delete(id string) (string, error)
}

// ChatEntity ....
type ChatEntity struct {
	ID        string `json:"id" bson:"_id"`
	RoomID    string `json:"room_id,omitempty" bson:"room_id,omitempty"`
	Message   string `json:"message" bson:"message"`
	Payload   string `json:"payload" bson:"payload"`
	Type      string `json:"type,omitempty" bson:"type,omitempty"`
	UserID    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at" bson:"deleted_at"`
}

// chatModel ...
type chatModel struct {
	DB     *mongo.Client
	DBName string
}

// NewChatModel ...
func NewChatModel(db *mongo.Client, dbName string) IChat {
	return &chatModel{DB: db, DBName: dbName}
}

// FindAllByRoom ...
func (model chatModel) FindAllByRoom(roomID, message, lastID string, limit int) (res []ChatEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("chats")

	// Option to set limit, offset and sorting
	l := int64(limit)
	findOptions := options.Find()
	findOptions.SetLimit(l)
	findOptions.SetSort(bson.D{{"created_at", -1}})

	match := []interface{}{
		bson.D{{"room_id", roomID}},
		bson.D{
			{"$or", []interface{}{
				bson.D{{"deleted_at", nil}},
				bson.D{{"deleted_at", ""}},
			}},
		},
	}

	if lastID != "" {
		match = append(match, bson.D{{"_id", bson.M{"$lt": lastID}}})
	}
	if message != "" {
		match = append(match, bson.D{{"message", primitive.Regex{Pattern: message, Options: "i"}}})
	}

	data, err := collection.Find(context.TODO(), bson.D{
		{
			"$and", match,
		},
	}, findOptions)
	if err != nil {
		return res, err
	}

	defer data.Close(context.TODO())
	for data.Next(context.TODO()) {
		var r ChatEntity
		data.Decode(&r)
		res = append(res, r)
	}
	if err := data.Err(); err != nil {
		return res, err
	}

	return res, err
}

// FindLast ...
func (model chatModel) FindLast(roomID string) (res ChatEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("chats")

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	err = collection.FindOne(context.TODO(), bson.D{
		{
			"$and", []interface{}{
				bson.D{{"room_id", roomID}},
				bson.D{
					{"$or", []interface{}{
						bson.D{{"deleted_at", nil}},
						bson.D{{"deleted_at", ""}},
					}},
				},
			},
		}}, findOptions).Decode(&res)

	return res, err
}

// Store ...
func (model chatModel) Store(body *ChatEntity) (res string, err error) {
	// Define ID
	body.ID = primitive.NewObjectIDFromTimestamp(time.Now().UTC()).Hex()

	collection := model.DB.Database(model.DBName).Collection("chats")
	_, err = collection.InsertOne(context.TODO(), body)
	if err != nil {
		return res, err
	}

	res = body.ID

	return res, err
}

// Delete ...
func (model chatModel) Delete(id string) (res string, err error) {
	collection := model.DB.Database(model.DBName).Collection("chats")

	filter := bson.D{{"_id", id}}
	now := time.Now().UTC()
	update := bson.M{"$set": bson.M{
		"updated_at": now.Format(time.RFC3339), "deleted_at": now.Format(time.RFC3339),
	}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return res, err
	}

	res = id

	return res, err
}
