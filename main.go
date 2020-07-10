package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
	"github.com/gofiber/limiter"
	"github.com/gofiber/logger"
	"github.com/joho/godotenv"
	"github.com/rbo13/chibyurl/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Environment File cannot be loaded: %v", err)
		return
	}

	var PORT = os.Getenv("PORT")

	db := dbConnect("chiby")
	if db == nil {
		log.Fatalf("Cannot connect to database!")
		return
	}
	collection = db.Collection("urls")

	server := fiber.New()

	// middlewares
	server.Use(
		middleware.Recover(),
		logger.New(),
		helmet.New(),
		// 5 requests per 15 seconds
		limiter.New(limiter.Config{
			Timeout: 15,
			Max:     5,
		}),
	)

	// serve static files
	server.Static("/", "./public")

	server.Get("/generate", func(ctx *fiber.Ctx) {
		url := new(model.URL)

		ctx.Status(http.StatusOK).JSON(fiber.Map{
			"data": url.Generate(),
		})
	})

	server.Get("/", func(ctx *fiber.Ctx) {
		urls := []model.URL{}
		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"code":    http.StatusNotFound,
				"message": err.Error(),
				"data":    nil,
			})
		}

		for cursor.Next(context.TODO()) {
			var url model.URL

			if err := cursor.Decode(&url); err != nil {
				ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"code":    http.StatusBadRequest,
					"message": err.Error(),
					"data":    nil,
				})
				return
			}
			urls = append(urls, url)
		}

		ctx.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
			"code":    http.StatusOK,
			"message": "URLS Retrieved",
			"data":    urls,
		})
	})

	server.Post("/", func(ctx *fiber.Ctx) {
		url := new(model.URL)

		if err := ctx.BodyParser(url); err != nil {
			ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		// check if alias is not blank, if blank create
		// a default alias.
		if url.Alias == "" {
			url.Alias = url.Generate()
		}

		// create index
		_, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.M{
					"alias": 1,
				},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			log.Fatalf("Cannot create document index: %v", err)
			return
		}

		result, err := collection.InsertOne(context.TODO(), url)
		if err != nil {
			ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		ctx.Status(http.StatusCreated).JSON(fiber.Map{
			"success": true,
			"code":    http.StatusCreated,
			"message": "Success!",
			"payload": url,
			"data":    result,
		})
	})

	server.Listen(PORT)
}

func dbConnect(dbName string) *mongo.Database {
	connString := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(connString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Cannot create MongoDB Client: %v", err)
		return nil
	}

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//To close the connection at the end
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
		return nil
	}

	log.Print("\n\n Successfully Connected to Mongo Database \n\n")
	return client.Database(dbName)
}
