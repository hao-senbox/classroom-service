package main

import (
	"classroom-service/config"
	"classroom-service/internal/room"
	"classroom-service/pkg/constants"
	"classroom-service/pkg/consul"
	"classroom-service/pkg/zap"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load env
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	consulConn := consul.NewConsulConn(logger, cfg)
	consulConn.Connect()

	// Handle OS signal để deregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server... De-registering from Consul...")
		consulConn.Deregister()
		os.Exit(0)
	}()
	
	c := cron.New(cron.WithSeconds())
	classroomCollection := mongoClient.Database(cfg.MongoDB).Collection("classroom")
	assginCollection := mongoClient.Database(cfg.MongoDB).Collection("assgin")
	systemConfig := mongoClient.Database(cfg.MongoDB).Collection("system_config")
	notification := mongoClient.Database(cfg.MongoDB).Collection("notification")
	classroomRepository := room.NewRoomRepository(classroomCollection, assginCollection, systemConfig, notification)
	classroomService := room.NewRoomService(classroomRepository)
	classroomHandler := room.NewRoomHandler(classroomService)

	r := gin.Default()

	room.RegisterRoutes(r, classroomHandler)
	_, err = c.AddFunc("0 */1 * * * *", func() {
		log.Println("🔄 Cron master running...")
		ctx := context.WithValue(context.Background(), constants.TokenKey, os.Getenv("CRON_SERVICE_TOKEN"))
		if err := classroomService.CronNotifications(ctx); err != nil {
			log.Printf("CronEventNotifications failed: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("AddFunc error: %v", err)
	}

	c.Start()
	defer c.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}
