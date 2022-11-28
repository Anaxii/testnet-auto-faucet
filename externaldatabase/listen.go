package externaldatabase

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListenForChanges(uri string, accountChannel chan string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	requestsCollection := client.Database("PuffinTestnet").Collection("approved")

	pipeline := bson.D{{"$match", bson.D{{"operationType", "insert"}}}}

	stream, err := requestsCollection.Watch(context.TODO(), mongo.Pipeline{pipeline})
	if err != nil {
		panic(err)
	}

	defer stream.Close(context.TODO())

	log.Info("Database listener running")

	for stream.Next(context.TODO()) {
		var newInsert primitive.D
		if err := stream.Decode(&newInsert); err != nil {
			log.Warn(err)
			continue
		}

		_data, err := json.Marshal(newInsert.Map()["fullDocument"])
		if err != nil {
			log.Warn(err)
			continue
		}

		var data []map[string]interface{}
		err = json.Unmarshal(_data,  &data)

		if err != nil {
			log.Warn(err)
			continue
		}

		for  _, v := range data {
			if v["Key"] == "wallet_address" {
				accountChannel <- fmt.Sprintf("%v", v["Value"])
			}
		}
	}
}
