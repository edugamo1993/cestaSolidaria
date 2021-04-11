package shopping

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-solidary/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

const (
	ShoppingCollection = "shoppingsystem"
)

var ErrSystemNotFound = errors.New("shopping system not found")
var ErrSystemAlreadyExist = errors.New("shopping system already exist")

type ShoppingSystem struct {
	ShoppingName string  `json:"shoppingName"`                   // Name that business choose for the shopping system
	BusinessCIF  string  `json:"businessCIF" bson:"businessCIF"` // BusinessCIF related to this system
	Accumulated  float64 `json:"accumulated" bson:"accumulated"` // Shopping ticket accumulated amount in €
	BasketValue  float64 `json:"basketValue" bson:"basketValue"` // Shopping ticket price in €
}

func (b *ShoppingSystem) AddMoney(c *config.Config, amount float64) error {
	b, err := GetShoppingSystem(c, b.BusinessCIF)
	if err != nil {
		return err
	}
	//Filter
	field, ok := reflect.TypeOf(b).Elem().FieldByName("BusinessCIF")
	if !ok {
		return fmt.Errorf("Error reflecting field")
	}
	bfilter := bson.M{string(field.Tag): bson.M{"$eq": b.BusinessCIF}}

	//Append amount
	b.Accumulated += amount
	//Serialice Data
	bsonData := bson.D{}
	bmarshal, err := bson.Marshal(b)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bmarshal, &bsonData)
	if err != nil {
		return err
	}
	err = c.Mongo.UpdateData(ShoppingCollection, bfilter, bsonData)
	return err
}

func GetShoppingSystem(c *config.Config, businessCIF string) (b *ShoppingSystem, err error) {

	res, err := c.Mongo.FindOne(bson.DocElem{Name: "businessCIF", Value: businessCIF}, ShoppingCollection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSystemNotFound
		}
		return nil, err

	}

	data, err := json.Marshal(res)
	err = json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return
}

func CreateShoppingSystem(c *config.Config, businessCIF string, shoppingValue float64, shoppingName string) (b *ShoppingSystem, err error) {
	_, err = GetShoppingSystem(c, businessCIF)
	if err != nil && err != ErrSystemNotFound {
		log.Println("error getting shopping system: ", err)
		return nil, err
	}
	if err == nil {
		log.Println("error creating shopping system: already exist")
		return nil, ErrSystemAlreadyExist
	}
	b = &ShoppingSystem{
		ShoppingName: shoppingName,
		BusinessCIF:  businessCIF,
		BasketValue:  shoppingValue,
		Accumulated:  0,
	}
	bmarshal, err := bson.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("error marshaling bson : %s", err)
	}
	bsonData := bson.DocElem{}
	err = bson.Unmarshal(bmarshal, &bsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bson : %s", err)
	}

	_, err = c.Mongo.InsertData(ShoppingCollection, bsonData)

	if err != nil {
		return nil, fmt.Errorf("error inserting basketSystem: %s", err)
	}

	return b, nil
}
