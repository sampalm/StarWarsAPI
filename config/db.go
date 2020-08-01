package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sampalm/starwars/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbuser = os.Getenv("MONGO_USER")
var dbpass = os.Getenv("MONGO_PASS")

func Conn() {
	// DATABASE CONNECTION
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.r7mjg.mongodb.net/<dbname>?retryWrites=true&w=majority", dbuser, dbpass)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer cancel()
	// DATABASE COLLETIONS
	db := client.Database("starwars")
	controllers.PlanetaCollection(db)
}
