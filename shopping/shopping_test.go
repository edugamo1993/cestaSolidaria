package shopping

import (
	"go-solidary/config"
	"go-solidary/mongo"
	"testing"
)

var c = &config.Config{
	Server: config.ServerConfig{},
	Mongo: mongo.Mongo{
		Addr:     "localhost:27027",
		DB:       "gosolidary",
		User:     "user1",
		Password: "example",
	},
}

func TestCreateShoppingSystem(t *testing.T) {

	res, err := CreateShoppingSystem(c, "123456-F", 30, "cestaSolidaria")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if res == nil {
		t.Log("error: res is nil")
		t.FailNow()
	}
}
