package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo type
type Mongo struct {
	Addr     string `json:"addr"`
	DB       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//GetAddr returns Addr property
func (m *Mongo) GetAddr() string {
	return m.Addr
}

//GetDatabase returns Database property
func (m *Mongo) GetDatabase() string {
	return m.DB
}

//GetUser returns User property
func (m *Mongo) GetUser() string {
	return m.User
}

//GetPassword returns Password property
func (m *Mongo) GetPassword() string {
	return m.Password
}

//NewSession returns MongoConnection
func (m *Mongo) NewSession(coll string) (ctx context.Context, collection *mongo.Collection, err error) {
	server := m.GetAddr()
	user := m.GetUser()
	password := m.GetPassword()
	addrToConnect := fmt.Sprintf("mongodb://%s:%s@%s/%s", user, password, server, m.GetDatabase())
	fmt.Println(addrToConnect)
	client, err := mongo.NewClient(options.Client().ApplyURI(addrToConnect))

	if err != nil {
		return nil, nil, err
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	collection = client.Database(m.GetDatabase()).Collection(coll)
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}
	return
}

//InsertData insert data in a table
func (m *Mongo) InsertData(collection string, data interface{}) (idInserted interface{}, err error) {
	ctx, coll, err := m.NewSession(collection)
	if err != nil {
		return nil, err
	}
	res, err := coll.InsertOne(ctx, data)
	if err != nil {
		if err.(mongo.WriteException).WriteErrors != nil && err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
			return nil, fmt.Errorf("duplicated")
		}
		return
	}

	idInserted = res.InsertedID
	fmt.Println(idInserted)
	return
}

func (m *Mongo) UpdateData(collection string, filter, data interface{}) error {
	ctx, coll, err := m.NewSession(collection)
	if err != nil {
		return err
	}
	res, err := coll.UpdateOne(ctx, filter, data)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("modifiedCount is %d", res.ModifiedCount)
	}

	return nil
}

//FindAll insert data in a table
func (m *Mongo) FindAll(query interface{}, collection string) (result []interface{}, err error) {
	ctx, coll, err := m.NewSession(collection)
	if err != nil {
		return
	}

	cur, err := coll.Find(ctx, query)
	if err != nil {
		log.Printf("error doing find request: %s", err)
		return nil, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var r interface{}
		err := cur.Decode(&r)
		if err != nil {
			log.Fatalf("error decoding: %s", err)
		}
		result = append(result, r)
	}
	if len(result) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return
}

//FindOne insert data in a table
func (m *Mongo) FindOne(query interface{}, collection string) (result interface{}, err error) {
	ctx, coll, err := m.NewSession(collection)
	if err != nil {
		return nil, err
	}
	err = coll.FindOne(ctx, query).Decode(&result)
	if err != nil {
		log.Println("findOne:", err)
		return nil, err
	}
	return
}

//DeleteOne insert data in a table
func (m *Mongo) DeleteOne(query interface{}, collection string) (err error) {
	ctx, coll, err := m.NewSession(collection)
	if err != nil {
		return err
	}

	res, err := coll.DeleteOne(ctx, query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if res.DeletedCount != 1 {
		return fmt.Errorf("deletedCount is %d", res.DeletedCount)
	}
	return
}
