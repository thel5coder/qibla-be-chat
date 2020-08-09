package mongomodel

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// RoomTypePrivate ...
	RoomTypePrivate = "private"
	// RoomTypeGroup ...
	RoomTypeGroup = "group"
)

// IRoom ...
type IRoom interface {
	FindAllByParticipant(participantID, message, lastID string, limit int) ([]RoomEntity, error)
	FindByID(id string) (RoomEntity, error)
	FindByIDParticipant(id, participantID string) (RoomEntity, error)
	FindByProfilePicture(userID, profilePicture string) (RoomEntity, error)
	FindPrivateByUser(userID, userParticipantID string) (RoomEntity, error)
	Store(body *RoomEntity) (string, error)
	Update(body *RoomEntity) (string, error)
	Delete(id string) (string, error)
}

// RoomEntity ....
type RoomEntity struct {
	ID                string `json:"id" bson:"_id"`
	Type              string `json:"type" bson:"type"`
	Name              string `json:"name" bson:"name"`
	ProfilePicture    string `json:"profile_picture" bson:"profile_picture"`
	Description       string `json:"description" bson:"description"`
	UserID            string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserParticipantID string `json:"user_participant_id,omitempty" bson:"user_participant_id,omitempty"`
	CreatedAt         string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt         string `json:"deleted_at" bson:"deleted_at"`
}

// roomModel ...
type roomModel struct {
	DB     *mongo.Client
	DBName string
}

// NewRoomModel ...
func NewRoomModel(db *mongo.Client, dbName string) IRoom {
	return &roomModel{DB: db, DBName: dbName}
}

// FindAllByMessage ...
func (model roomModel) FindAllByParticipant(participantID, message, lastID string, limit int) (res []RoomEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")

	// Option to set limit, offset and sorting
	l := int64(limit)

	// Match interface
	match := []interface{}{
		bson.D{{"participant.user_id", participantID}},
		bson.D{
			{"$or", []interface{}{
				bson.D{{"participant.deleted_at", nil}},
				bson.D{{"participant.deleted_at", ""}},
			}},
		},
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
		match = append(match, bson.D{
			{"$or", []interface{}{
				bson.D{{"name", primitive.Regex{Pattern: message, Options: "i"}}},
				bson.D{{"message", primitive.Regex{Pattern: message, Options: "i"}}},
			}},
		})
	}

	query := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "participants",
				"localField":   "_id",
				"foreignField": "room_id",
				"as":           "participant",
			},
		},
		{"$unwind": "$participant"},
		{
			"$match": bson.D{
				{
					"$and", match,
				},
			},
		},
		{
			"$sort": bson.M{"updated_at": -1},
		},
		{
			"$limit": l,
		},
	}

	data, err := collection.Aggregate(context.TODO(), query)
	if err != nil {
		return res, err
	}

	defer data.Close(context.TODO())
	for data.Next(context.TODO()) {
		var r RoomEntity
		data.Decode(&r)
		res = append(res, r)
	}
	if err := data.Err(); err != nil {
		return res, err
	}

	return res, err
}

// FindByID ...
func (model roomModel) FindByID(id string) (res RoomEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	err = collection.FindOne(context.TODO(), bson.D{
		{
			"$and", []interface{}{
				bson.D{{"_id", id}},
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

// FindByIDParticipant ...
func (model roomModel) FindByIDParticipant(id, participantID string) (res RoomEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")

	// Match interface
	match := []interface{}{
		bson.D{{"_id", id}},
		bson.D{{"participant.user_id", participantID}},
		bson.D{
			{"$or", []interface{}{
				bson.D{{"participant.deleted_at", nil}},
				bson.D{{"participant.deleted_at", ""}},
			}},
		},
		bson.D{
			{"$or", []interface{}{
				bson.D{{"deleted_at", nil}},
				bson.D{{"deleted_at", ""}},
			}},
		},
	}

	query := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "participants",
				"localField":   "_id",
				"foreignField": "room_id",
				"as":           "participant",
			},
		},
		{"$unwind": "$participant"},
		{
			"$match": bson.D{
				{
					"$and", match,
				},
			},
		},
		{
			"$sort": bson.M{"updated_at": -1},
		},
		{
			"$limit": 1,
		},
	}

	data, err := collection.Aggregate(context.TODO(), query)
	if err != nil {
		return res, err
	}

	defer data.Close(context.TODO())
	for data.Next(context.TODO()) {
		data.Decode(&res)
	}
	if err := data.Err(); err != nil {
		return res, err
	}

	return res, nil
}

// FindByProfilePicture ...
func (model roomModel) FindByProfilePicture(userID, profilePicture string) (res RoomEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	err = collection.FindOne(context.TODO(), bson.D{
		{
			"$and", []interface{}{
				bson.D{{"user_id", userID}},
				bson.D{{"profile_picture", profilePicture}},
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

// FindPrivateByUser ...
func (model roomModel) FindPrivateByUser(userID, userParticipantID string) (res RoomEntity, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	err = collection.FindOne(context.TODO(), bson.D{
		{
			"$and", []interface{}{
				bson.D{{
					"$or", []interface{}{
						bson.D{{
							"$and", []interface{}{
								bson.D{{"user_id", userID}},
								bson.D{{"user_participant_id", userParticipantID}},
							},
						}},
						bson.D{{
							"$and", []interface{}{
								bson.D{{"user_id", userParticipantID}},
								bson.D{{"user_participant_id", userID}},
							},
						}},
					}}},
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
func (model roomModel) Store(body *RoomEntity) (res string, err error) {
	// Define ID
	body.ID = primitive.NewObjectIDFromTimestamp(time.Now().UTC()).Hex()

	collection := model.DB.Database(model.DBName).Collection("rooms")
	_, err = collection.InsertOne(context.TODO(), body)
	if err != nil {
		return res, err
	}

	res = body.ID

	return res, err
}

// Update ...
func (model roomModel) Update(body *RoomEntity) (res string, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")
	filter := bson.D{{
		"$and", []interface{}{
			bson.D{{"_id", body.ID}},
			bson.D{{"type", RoomTypeGroup}},
			bson.D{{"user_id", body.UserID}},
			bson.D{{
				"$or", []interface{}{
					bson.D{{"deleted_at", nil}},
					bson.D{{"deleted_at", ""}},
				}},
			}},
	}}
	update := bson.M{"$set": bson.M{
		"name": body.Name, "profile_picture": body.ProfilePicture,
		"description": body.Description, "updated_at": body.UpdatedAt,
	}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return res, err
	}

	res = body.ID

	return res, err
}

// Delete ...
func (model roomModel) Delete(id string) (res string, err error) {
	collection := model.DB.Database(model.DBName).Collection("rooms")
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
