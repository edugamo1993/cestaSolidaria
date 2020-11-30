package mongo

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var mongoConfig = Mongo{
	Addr:     "localhost:27027",
	DB:       "gosolidary",
	User:     "user1",
	Password: "example",
}

func TestFindAll(t *testing.T) {
	t.Log("a")

	res, err := mongoConfig.FindAll(struct{}{}, "test")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	fmt.Println("result: ", res)
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
		fmt.Println(err.Error())
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
