package business

import (
	"encoding/json"
	"fmt"
	"go-solidary/config"

	"go.mongodb.org/mongo-driver/bson"
)

//Business type
type Business struct {
	ID         string `json:"-" bson:"_id"`
	CIF        string `json:"cif" bson:"cif"`
	CommonName string `json:"commonName" bson:"commonName"`
	OwnerName  string `json:"ownerName" bson:"ownerName"`
	Phone      string `json:"phone" bson:"phone"`
	Email      string `json:"email" bson:"email"`
	Verified   bool   `json:"-" bson:"verified"`
}

const (
	businessTable = "business"
)

//InsertBusiness function create business on database
func InsertBusiness(c *config.Config, data []byte) (*Business, error) {
	bjsonData := bson.D{}
	b := Business{}
	err := json.Unmarshal(data, &b)
	if err != nil {
		return nil, err
	}
	b.ID = b.CIF
	bmarshal, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(bmarshal, &bjsonData)
	if err != nil {
		return nil, err
	}
	id, err := c.Mongo.InsertData(businessTable, bjsonData)
	if err != nil {
		return nil, err
	}
	fmt.Println(id)
	return &b, err
}

//GetBusinessBy function returns []Business filtered by any mongo field
func GetBusinessBy(c *config.Config, field, value string) (b []Business, err error) {

	s := fmt.Sprintf("{%s: %s}", field, value)
	result, err := c.Mongo.FindAll(s, businessTable)
	if err != nil {
		return nil, err
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonResult, &b)
	if err != nil {
		return nil, err

	}

	return
}
