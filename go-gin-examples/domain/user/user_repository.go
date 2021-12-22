package user

import (
	"cathub.me/go-web-examples/pkg/data"
	"cathub.me/go-web-examples/pkg/database"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

// _userRepository Singleton
var _userRepository UserRepository
var _onceUserRepository sync.Once

func GetUserRepository() UserRepository {
	_onceUserRepository.Do(func() {
		_userRepository = &userRepository{collection: database.GetMongoDatabase().Collection("user")}
	})
	return _userRepository
}

type UserRepository interface {
	Collection() *mongo.Collection
	Insert(ctx context.Context, document *User) error
	InsertAll(ctx context.Context, documents []*User) error

	DeleteAll(ctx context.Context) error
	DeleteById(ctx context.Context, id primitive.ObjectID) error
	DeleteByIds(ctx context.Context, ids []primitive.ObjectID) error

	UpdateById(ctx context.Context, document *User) (*User, error)

	FindAll(ctx context.Context, pageable data.Pageable) (*data.Pageable, []*User, error)
	FindOne(ctx context.Context, filter interface{}) (*User, error)
	FindById(ctx context.Context, id primitive.ObjectID) (*User, error)
	Count(ctx context.Context) (int64, error)

	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)

	FindAllByActivatedIsFalseAndCreatedDateLessThan(ctx context.Context, time time.Time, request data.Pageable) ([]*User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func (r *userRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *userRepository) Insert(ctx context.Context, document *User) error {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	if _, err := r.collection.InsertOne(ctxWt, document); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) InsertAll(ctx context.Context, documents []*User) error {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	_documents := make([]interface{}, len(documents))
	for i, document := range documents {
		_documents[i] = document
	}

	if _, err := r.collection.InsertMany(ctxWt, _documents); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteAll(ctx context.Context) error {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	if _, err := r.collection.DeleteMany(ctxWt, bson.M{}); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	if _, err := r.collection.DeleteOne(ctxWt, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteByIds(ctx context.Context, ids []primitive.ObjectID) error {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	_ids := bson.D{}
	for _, id := range ids {
		_ids = append(_ids, bson.E{Key: "_id", Value: id})
	}

	if _, err := r.collection.DeleteMany(ctxWt, _ids); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateById(ctx context.Context, document *User) (*User, error) {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	updateAfterReturn := options.After
	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &updateAfterReturn,
	}

	var _document *User
	_ = r.collection.FindOneAndUpdate(ctxWt, bson.M{"_id": document.Id}, bson.M{"$set": document}, updateOptions).Decode(&_document)

	return _document, nil
}

func (r *userRepository) FindAll(ctx context.Context, pageable data.Pageable) (*data.Pageable, []*User, error) {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	filter := bson.M{}
	cursor, err := r.collection.Find(ctxWt, filter, database.GetFindOptions(pageable))
	if err != nil {
		return nil, nil, err
	}

	var documents []*User
	for cursor.Next(ctxWt) {
		var entity *User
		if err := cursor.Decode(&entity); err == nil {
			documents = append(documents, entity)
		}
	}

	count, err := r.collection.CountDocuments(ctxWt, filter)
	if err != nil {
		return nil, nil, err
	}

	result := pageable.Copy(count)
	return result, documents, nil
}

func (r *userRepository) FindById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	return r.FindOne(ctx, bson.M{"_id": id})
}

func (r *userRepository) FindOne(ctx context.Context, filter interface{}) (*User, error) {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	var document User
	result := r.collection.FindOne(ctxWt, filter)
	if err := result.Decode(&document); err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	count, err := r.collection.CountDocuments(ctxWt, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	return r.FindOne(ctx, bson.M{"username": username})
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return r.FindOne(ctx, bson.M{"email": email})
}

func (r *userRepository) FindAllByActivatedIsFalseAndCreatedDateLessThan(ctx context.Context, lteTime time.Time, request data.Pageable) ([]*User, error) {
	ctxWt, cancelFunc := context.WithTimeout(ctx, 60*time.Second)
	defer cancelFunc()

	filter := bson.M{"activated": false, "created_date": bson.M{"$lte": lteTime}}
	cursor, err := r.collection.Find(ctxWt, filter, database.GetFindOptions(request))
	if err != nil {
		return nil, err
	}

	var documents []*User
	for cursor.Next(ctxWt) {
		var entity *User
		if err := cursor.Decode(&entity); err == nil {
			documents = append(documents, entity)
		}
	}

	return documents, nil
}
