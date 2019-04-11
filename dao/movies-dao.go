package dao

import (
	"context"
	"fmt"
	"os"

	"github.com/zacsketches/crud/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MoviesDAO contains the connection info for the database
type MoviesDAO struct {
	User     string
	Server   string
	Database string
	client   *mongo.Client
}

const (
	COLLECTION = "movies"
)

// Connect establishes a persistent connection to the database or returns
// an error if it cannot connect.  Expects the user password to be stored
// in an environment variable entitled MONGO_PW.
func (m *MoviesDAO) Connect() error {
	ctx := context.TODO()

	pw, ok := os.LookupEnv("MONGO_PW")
	if !ok {
		return fmt.Errorf("unable to find MONGO_PW in the environment")
	}

	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s", m.User, pw, m.Server)

	// Set client options and connect
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	m.client = client

	return nil
}

// FindAll returns a slice of Movie
func (m *MoviesDAO) FindAll() ([]models.Movie, error) {
	ctx := context.TODO()
	collection := m.client.Database(m.Database).Collection(COLLECTION)
	var movies []models.Movie

	findOptions := options.Find()
	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem models.Movie
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		movies = append(movies, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

// FindByID returns a single movie identified by its id
func (m *MoviesDAO) FindByID(id string) (models.Movie, error) {
	ctx := context.TODO()
	var result models.Movie
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Movie{}, err
	}

	filter := bson.D{{"_id", objID}}

	err = m.client.Database(m.Database).Collection(COLLECTION).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.Movie{}, err
	}

	return result, nil
}

// Insert places another Movie document into the collection
func (m *MoviesDAO) Insert(movie models.Movie) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	bson := bson.M{"name": movie.Name, "description": movie.Description, "cover_image": movie.CoverImage}
	res, err := m.client.Database(m.Database).Collection(COLLECTION).InsertOne(ctx, bson)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete removes a document from the collection
func (m *MoviesDAO) Delete(movie models.Movie) (*mongo.DeleteResult, error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", movie.ID}}
	res, err := m.client.Database(m.Database).Collection(COLLECTION).DeleteOne(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MoviesDAO) Update(movie models.Movie) (*mongo.UpdateResult, error) {
	ctx := context.TODO()

	filter := bson.D{{"_id", movie.ID}}
	update := bson.D{{"$set",
		bson.D{
			{"name", movie.Name},
			{"cover_image", movie.CoverImage},
			{"description", movie.Description},
		},
	}}

	res, err := m.client.Database(m.Database).Collection(COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}
