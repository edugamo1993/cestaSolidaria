package business

import (
	"encoding/json"
	"fmt"
	"go-solidary/config"
	"log"
	"strconv"
	"strings"

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
	State      string `json:"state" bson:"state"`
	Address    string `json:"address" bson:"address"`
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

func validateCIF(cif string) bool {
	ctrlDigits := []string{"J", "A", "B", "C", "D", "E", "F", "G", "H", "I"}
	firstDigit := cif[0:1]
	centralDigits := cif[1 : len(cif)-1]
	ctrlDigit := cif[len(cif)-1:]

	a := 0
	for i := 1; i < len(centralDigits); i += 2 {
		digit, err := strconv.Atoi(string(centralDigits[i]))
		if err != nil {
			log.Printf("error: digit %d cannot be integer", digit)
		}
		a += digit
	}
	b := 0
	for i := 0; i < len(centralDigits); i += 2 {
		digit, err := strconv.Atoi(string(centralDigits[i]))
		if err != nil {
			log.Printf("error: digit %d cannot be integer", digit)
		}
		double := digit * 2
		var sumDigits int
		sumDigits = double/10 + double%10
		b += sumDigits
	}
	c := (a + b) % 10 //Last digit of (a+b)
	d := 10 - c       // 10 - C == ctrlDigit if number, == ctrlDigits[d] if string
	if strings.Contains("abeh", strings.ToLower(firstDigit)) {
		if strconv.Itoa(d) != ctrlDigit {
			return false
		}
	} else {
		if ctrlDigit != ctrlDigits[d] {
			return false
		}
	}
	return true
}
