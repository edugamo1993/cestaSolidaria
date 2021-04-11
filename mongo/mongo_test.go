package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var mongoConfig = Mongo{
	Addr:     "localhost:27027",
	DB:       "gosolidary",
	User:     "user1",
	Password: "example",
}

func TestMongo_FindAllNotFound(t *testing.T) {
	res, err := mongoConfig.FindAll(bson.DocElem{"_id", "abcdefg"}, "test")
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err.Error())
		t.FailNow()
	}
	if res != nil {
		t.FailNow()
	}
}

func TestMongo_FindOneNotFound(t *testing.T) {
	res, err := mongoConfig.FindOne(struct{}{}, "test")
	if err != nil && err != mongo.ErrNoDocuments {
		t.Log(err.Error())
		t.FailNow()
	}
	if res != nil {
		t.FailNow()
	}
}

func TestInsertData(t *testing.T) {

	data := struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}{
		"1234", "Edu",
	}

	b, err := bson.Marshal(data)
	res, err := mongoConfig.InsertData("test", b)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if res == nil {
		t.FailNow()
	}
	fmt.Println(res)
}

func TestUpdate(t *testing.T) {

	data := bson.M{"$set": bson.M{"_id": "1234", "name": "Eduardo"}}

	filter := bson.M{"_id": bson.M{"$eq": "1234"}}
	err := mongoConfig.UpdateData("test", filter, data)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

}

func TestDelete(t *testing.T) {
	filter := bson.M{"_id": bson.M{"$eq": "1234"}}
	err := mongoConfig.DeleteOne(filter, "test")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}
