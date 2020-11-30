package basket

import (
	"encoding/json"
	"fmt"
	"go-solidary/config"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

const (
	BasketCollection = "basketsystem"
)

type BasketSystem struct {
	BusinessCIF string `json:"businessCIF" bson:"businessCIF"`
	Accumulated int    `json:"accumulated" bson:"accumulated"`
	BasketValue int    `json:"basketValue" bson:"basketValue"`
}

func (b *BasketSystem) AddMoney(c *config.Config, ammount int) error {
	b, err := GetBasketSystem(c, b.BusinessCIF)
	if err != nil {
		return err
	}
	//Filter
	field, ok := reflect.TypeOf(b).Elem().FieldByName("BusinessCIF")
	if !ok {
		return fmt.Errorf("Error reflecting field")
	}
	bfilter := bson.M{string(field.Tag): bson.M{"$eq": b.BusinessCIF}}

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
	err = c.Mongo.UpdateData(BasketCollection, bfilter, bsonData)
	return err
}

func GetBasketSystem(c *config.Config, businessCIF string) (b *BasketSystem, err error) {

	res, err := c.Mongo.FindOne(`{businessCIF:"`+businessCIF+`"}`, BasketCollection)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(res)
	err = json.Unmarshal(data, &b)
	if err != nil {
		return nil, err
	}
	return
}

func CreateBasketSystem(c *config.Config, businessCIF string, basketValue int) (b *BasketSystem, err error) {
	b = &BasketSystem{
		BusinessCIF: businessCIF,
		BasketValue: b.BasketValue,
		Accumulated: 0,
	}
	bmarshal, err := bson.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("error marshaling bson : %s", err)
	}
	bsonData := bson.D{}
	err = bson.Unmarshal(bmarshal, &bsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bson : %s", err)
	}

	_, err = c.Mongo.InsertData(BasketCollection, bsonData)

	if err != nil {
		return nil, fmt.Errorf("error inserting basketSystem: %s", err)
	}

	return b, nil
}
