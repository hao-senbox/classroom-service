package main

import (
	"classroom-service/config"
	"classroom-service/internal/assign"
	"classroom-service/internal/classroom"
	"classroom-service/internal/language"
	"classroom-service/internal/leader"
	"classroom-service/internal/region"
	"classroom-service/internal/room"
	"classroom-service/internal/term"
	"classroom-service/internal/user"
	"classroom-service/pkg/consul"
	"classroom-service/pkg/zap"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	// c := cron.New(cron.WithSeconds())
	// assginCollection := mongoClient.Database(cfg.MongoDB).Collection("assgin")
	// systemConfig := mongoClient.Database(cfg.MongoDB).Collection("system_config")
	// notification := mongoClient.Database(cfg.MongoDB).Collection("notification")
	// classCollection := mongoClient.Database(cfg.MongoDB).Collection("classroom")
	// leader := mongoClient.Database(cfg.MongoDB).Collection("leader")

	roomService := room.NewRoomService(consulClient)
	userService := user.NewUserService(consulClient)
	languageService := language.NewUserService(consulClient)
	termService := term.NewTermService(consulClient)

	regionCollection := mongoClient.Database(cfg.MongoDB).Collection("region")
	classroomCollection := mongoClient.Database(cfg.MongoDB).Collection("classroom")
	assignCollection := mongoClient.Database(cfg.MongoDB).Collection("assign")
	leaderCollection := mongoClient.Database(cfg.MongoDB).Collection("leader")
	assignTemplateCollection := mongoClient.Database(cfg.MongoDB).Collection("assign_template")
	leaderTemplateCollection := mongoClient.Database(cfg.MongoDB).Collection("leader_template")

	leaderRepository := leader.NewLeaderRepository(leaderCollection, leaderTemplateCollection)
	leaderService := leader.NewLeaderService(leaderRepository)
	leaderHandler := leader.NewLeaderHandler(leaderService)

	assignRepository := assign.NewAssignRepository(assignCollection, assignTemplateCollection)
	assignService := assign.NewAssignService(assignRepository)
	assignHandler := assign.NewAssignHandler(assignService)

	classroomRepository := classroom.NewClassroomRepository(classroomCollection)
	classroomService := classroom.NewClassroomService(classroomRepository, assignRepository, userService, leaderRepository, languageService, termService)
	classroomHandler := classroom.NewClassroomHandler(classroomService)

	regionRepository := region.NewRegionRepository(regionCollection)
	regionService := region.NewRegionService(regionRepository, classroomRepository, assignRepository, userService, roomService, leaderRepository, languageService)
	regionHandler := region.NewRegionHandler(regionService)

	// classroomRepository := class.NewClassRepository(assginCollection, systemConfig, notification, leader, classCollection)
	// classroomService := class.NewClassService(classroomRepository, roomService, userService)
	// classroomHandler := class.NewClassHandler(classroomService)

	r := gin.Default()

	leader.RegisterRoutes(r, leaderHandler)
	assign.RegisterRoutes(r, assignHandler)
	classroom.RegisterRoutes(r, classroomHandler)
	region.RegisterRoutes(r, regionHandler)

	// _, err = c.AddFunc("0 0 0 * * *", func() {
	// 	log.Println("🔄 Cron master running...")
	// 	ctx := context.WithValue(context.Background(), constants.TokenKey, os.Getenv("CRON_SERVICE_TOKEN"))
	// 	if err := classroomService.CronNotifications(ctx); err != nil {
	// 		log.Printf("CronEventNotifications failed: %v", err)
	// 	}
	// })
	// if err != nil {
	// 	log.Fatalf("AddFunc error: %v", err)
	// }

	// c.Start()
	// defer c.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped with error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
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
