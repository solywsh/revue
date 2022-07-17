package mongo_service

import (
	"context"
	"github.com/solywsh/qqBot-revue/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var (
	mongoOnce sync.Once
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo() (m *Mongo, f bool) {
	mongoOnce.Do(func() {
		yamlConf, err := conf.NewConf("./config.yaml")
		// Set DbClient options
		clientOptions := options.Client().ApplyURI(yamlConf.Database.Mongo.HImgDB.Url)
		// Connect to MongoDB
		m.Client, err = mongo.Connect(context.TODO(), clientOptions)
		// Check the connection
		err = m.Client.Ping(context.TODO(), nil)
		if err != nil {
			log.Println("NewMongo error:", err)
			f = false
			return
		} else {
			f = true
		}
		//DbCollection = DbClient.Database("revue").Collection("himg")
		//log.Println("Connected to MongoDB!")
	})
	return m, f
}
