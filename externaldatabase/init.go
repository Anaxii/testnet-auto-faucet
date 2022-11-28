package externaldatabase

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetAccounts(uri string) map[string]bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer client.Disconnect(ctx)

	requestsCollection := client.Database("PuffinTestnet").Collection("approved")

	cursor, err := requestsCollection.Find(context.TODO(), bson.D{{"status", "approved"}})
	var results []VerificationRequest
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Error("Could not get accounts from external database")
		return nil
	}

	accounts := map[string]bool{}

	for _, result := range results {
		accounts[result.WalletAddress] = true
	}

	return accounts

}
