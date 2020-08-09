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
	// ParticipantTypeAdmin ...
	ParticipantTypeAdmin = "admin"
	// ParticipantTypeUser ...
	ParticipantTypeUser = "user"
)

// IParticipant ...
type IParticipant interface {
	SelectAllByRoom(roomID string) ([]ParticipantEntity, error)
	FindAllByRoom(roomID, lastID string, limit int) ([]ParticipantEntity, error)
	FindByRoomParticipant(roomID, userID string) (ParticipantEntity, error)
	Store(body *ParticipantEntity) (string, error)
	Delete(id string) (string, error)
	DeleteByRoomParticipant(roomID, userID string) error
}

// ParticipantEntity ....
type ParticipantEntity struct {
	ID        string `json:"id" bson:"_id"`
	RoomID    string `json:"room_id,omitempty" bson:"room_id,omitempty"`
	UserID    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Type      string `json:"type,omitempty" bson:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at" bson:"deleted_at"`
}

// participantModel ...
type participantModel struct {
	DB     *mongo.Client
	DBName string
}

// NewParticipantModel ...
func NewParticipantModel(db *mongo.Client, dbName string) IParticipant {
	return &participantModel{DB: db, DBName: dbName}
}

// SelectAllByRoom ...
func (model participantModel) SelectAllByRoom(roomID string) (res []ParticipantEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("participants")

	// Option to set limit, offset and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"type", 1}})

	match := []interface{}{
		bson.D{{"room_id", roomID}},
		bson.D{
			{"$or", []interface{}{
				bson.D{{"deleted_at", nil}},
				bson.D{{"deleted_at", ""}},
			}},
		},
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
		var r ParticipantEntity
		data.Decode(&r)
		res = append(res, r)
	}
	if err := data.Err(); err != nil {
		return res, err
	}

	return res, err
}

// FindAllByRoom ...
func (model participantModel) FindAllByRoom(roomID, lastID string, limit int) (res []ParticipantEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("participants")

	// Option to set limit, offset and sorting
	l := int64(limit)
	findOptions := options.Find()
	findOptions.SetLimit(l)
	findOptions.SetSort(bson.D{{"type", 1}})

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
		var r ParticipantEntity
		data.Decode(&r)
		res = append(res, r)
	}
	if err := data.Err(); err != nil {
		return res, err
	}

	return res, err
}

// FindByRoomParticipant ...
func (model participantModel) FindByRoomParticipant(roomID, userID string) (res ParticipantEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("participants")

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	err = collection.FindOne(context.TODO(), bson.D{
		{
			"$and", []interface{}{
				bson.D{{"room_id", roomID}},
				bson.D{{"user_id", userID}},
				bson.D{{
					"$or", []interface{}{
						bson.D{{"deleted_at", nil}},
						bson.D{{"deleted_at", ""}},
					}},
				}},
		},
	}, findOptions).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Store ...
func (model participantModel) Store(body *ParticipantEntity) (res string, err error) {
	// Define ID
	body.ID = primitive.NewObjectIDFromTimestamp(time.Now().UTC()).Hex()

	collection := model.DB.Database(model.DBName).Collection("participants")
	_, err = collection.InsertOne(context.TODO(), body)
	if err != nil {
		return res, err
	}

	res = body.ID

	return res, err
}

// Delete ...
func (model participantModel) Delete(id string) (res string, err error) {
	collection := model.DB.Database(model.DBName).Collection("participants")

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

// DeleteByRoomParticipant ...
func (model participantModel) DeleteByRoomParticipant(roomID, userID string) (err error) {
	collection := model.DB.Database(model.DBName).Collection("participants")

	filter := bson.D{{
		"$and", []interface{}{
			bson.D{{"room_id", roomID}},
			bson.D{{"user_id", userID}},
			bson.D{{
				"$or", []interface{}{
					bson.D{{"deleted_at", nil}},
					bson.D{{"deleted_at", ""}},
				}},
			}},
	}}

	now := time.Now().UTC()
	update := bson.M{"$set": bson.M{
		"updated_at": now.Format(time.RFC3339), "deleted_at": now.Format(time.RFC3339),
	}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
}
