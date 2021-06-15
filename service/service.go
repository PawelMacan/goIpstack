package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"goipstack/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getMongoDbConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return client, nil
}

func getMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection()
	if err != nil {
		return nil, err
	}
	collection := client.Database(DbName).Collection(CollectionName)
	return collection, nil
}

func CreateIpGeoLocation(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(os.Getenv("DB_NAME"), os.Getenv("COLLECTION_NAME"))
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var ip string
	if c.Params("ip") != "" {
		ip = c.Params("ip")
	}
	jsonStr, err := http.Get(os.Getenv("IP_SERVICE_URL") + ip + "?access_key=" + os.Getenv("ACCESS_KEY"))
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	responseData, err := ioutil.ReadAll(jsonStr.Body)
	if err != nil {
		log.Fatal(err)
	}
	var geoLocation model.GeoLocation
	json.Unmarshal(responseData, &geoLocation)

	res, err := collection.InsertOne(context.Background(), geoLocation)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	c.JSON(fiber.Map{
		"id":       res.InsertedID.(primitive.ObjectID).Hex(),
		"geo_data": geoLocation,
	})
}

func GetIpGeoLocation(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(os.Getenv("DB_NAME"), os.Getenv("COLLECTION_NAME"))
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var filter bson.M = bson.M{}
	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}
	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	cur.All(context.Background(), &results)
	if results == nil {
		c.SendStatus(404)
		return
	}
	json, _ := json.Marshal(results)
	c.Send(json)
	return
}

func DeleteIpGeoLocation(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(os.Getenv("DB_NAME"), os.Getenv("COLLECTION_NAME"))
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var filter bson.M = bson.M{}
	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Delete removed record %v \n", c.Params("id"))
	c.Status(204).Send()
}
