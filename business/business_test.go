package business

import (
	"fmt"
	"go-solidary/config"
	"go-solidary/mongo"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var configMongo = mongo.Mongo{
	Addr:     "localhost:27027",
	DB:       "gosolidary",
	User:     "user1",
	Password: "example",
}

func deleteBusiness() error {
	err := configMongo.DeleteOne(bson.M{"cif": bson.M{"$eq": "my-cif"}}, "business")
	return err
}

func TestInsertBusiness_HappyPath(t *testing.T) {
	defer deleteBusiness()
	b := []byte(`{
		"cif":"my-cif",
		"commonName":"my-business",
		"ownerName":"ownerName",
		"phone":"21564869",
		"email":"email@email.com"
	}`)

	_, err := InsertBusiness(&config.Config{Mongo: configMongo}, b)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}
