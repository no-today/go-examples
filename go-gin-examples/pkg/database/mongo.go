package database

import (
	"cathub.me/go-web-examples/pkg/data"
	"cathub.me/go-web-examples/pkg/setting"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"sync"
	"time"
)

var _mongoDatabase *mongo.Database
var _onceMongoDatabase sync.Once

func GetMongoDatabase() *mongo.Database {
	_onceMongoDatabase.Do(func() {
		uri := setting.Mongodb.Uri

		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		clientOptions := options.Client().ApplyURI(uri)
		if setting.Mongodb.Debug {
			cmdMonitor := &event.CommandMonitor{
				Started: func(_ context.Context, evt *event.CommandStartedEvent) {
					if evt.CommandName == "find" || evt.CommandName == "aggregate" {
						log.Info().Str("module", "MongoDB").Int64("requestId", evt.RequestID).Str("commandName", evt.CommandName).Str("command", evt.Command.String()).Msg("")
					}
				},
				Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
				},
				Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
				},
			}
			clientOptions.SetMonitor(cmdMonitor)
		}

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal().Err(err).Msg("Connect mongodb failed")
		}

		cs, _ := connstring.ParseAndValidate(uri)
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatal().Err(err).Msg("Ping mongodb failed")
		}

		_mongoDatabase = client.Database(cs.Database)
	})
	return _mongoDatabase
}

func GetFindOptions(request data.Pageable) *options.FindOptions {
	findOptions := &options.FindOptions{}
	if request.Sorts != nil {
		sorts := bson.D{}
		for _, sort := range request.Sorts {
			sorts = append(sorts, bson.E{Key: sort.Field, Value: sort.Order})
		}
		findOptions.SetSort(sorts)
	}

	if request.Fields != nil {
		fields := bson.D{}
		for _, field := range request.Fields {
			fields = append(fields, bson.E{Key: field, Value: 1})
		}
		findOptions.SetProjection(fields)
	}

	findOptions.SetSkip((request.Page - 1) * request.Size)
	findOptions.SetLimit(request.Size)
	return findOptions
}
