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
	m         *Mongo
	f         bool
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo() (*Mongo, bool) {
	mongoOnce.Do(func() {
		var err error
		m = new(Mongo)
		yamlConf := conf.NewConf()
		// Set DbClient options
		clientOptions := options.Client().ApplyURI(yamlConf.Database.Mongo.HImgDB.Url)
		// Connect to MongoDB
		m.Client, err = mongo.Connect(context.TODO(), clientOptions)
		// Check the connection
		err = m.Client.Ping(context.TODO(), nil)
		if err != nil {
			log.Println("New Mongo Error: ", err)
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
