package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tosinjs/reminder-service/cmd/api/routes"
	firebaseinit "tosinjs/reminder-service/internal/firebase-init"
	mongoinit "tosinjs/reminder-service/internal/mongo-init"
	nmRepo "tosinjs/reminder-service/internal/repository/notificationRepo/mongoRepo"
	ntmRepo "tosinjs/reminder-service/internal/repository/notificationTokenRepo/mongoRepo"
	tmRepo "tosinjs/reminder-service/internal/repository/todoRepo/mongoRepo"
	"tosinjs/reminder-service/internal/service/authService"
	"tosinjs/reminder-service/internal/service/notificationService"
	"tosinjs/reminder-service/internal/service/notificationTokenService"
	"tosinjs/reminder-service/internal/service/reminderService"
	"tosinjs/reminder-service/internal/service/todoService"
	"tosinjs/reminder-service/internal/service/validationService"
	"tosinjs/reminder-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup() {
	config, err := utils.LoadConfig("./", "app", "env")
	if err != nil {
		log.Fatalf("Error GETIING CONFIG: %v", err)
		os.Exit(1)
	}

	//Database Setup
	mc, err := mongoinit.Init(context.Background(), config.MONGO_URI)
	if err != nil {
		fmt.Printf("mongo error: %v", err)
		os.Exit(1)
	}

	defer mongoinit.Disconnect(mc, context.Background())
	mongoDB := mc.Database(config.MONGODB)

	//Mongo Collections
	todoCollection := mongoDB.Collection("todo")
	notifTokenCollection := mongoDB.Collection("notificationTokens")
	notifCollection := mongoDB.Collection("notifications")
	//Notifications should expire after 1 week
	const hours_in_a_week = 24 * 7
	index := mongo.IndexModel{
		Keys: bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(
			int32((time.Hour * 2 * hours_in_a_week).Seconds()),
		),
	}
	_, err = notifCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		fmt.Printf("mongo index error: %v", err)
		os.Exit(1)
	}

	//Notification Tokens Live For 3 weeks
	index = mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(
			int32((time.Hour * 2 * hours_in_a_week).Seconds()),
		),
	}
	_, err = notifTokenCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		fmt.Printf("mongo index error: %v", err)
		os.Exit(1)
	}

	//Firebase Setup
	fcmClient, err := firebaseinit.Init(context.Background())
	if err != nil {
		fmt.Printf("error connecting to firebase: %v", err)
		os.Exit(1)
	}

	//GoCron Setup
	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	//Repo Setup
	notifRepo := nmRepo.New(notifCollection, context.Background())
	notifTokenRepo := ntmRepo.New(notifTokenCollection, context.Background())
	todoRepo := tmRepo.New(todoCollection, context.Background())

	//Service Setup
	authSVC := authService.New(config.JWTSECRET)
	validationSVC := validationService.New()
	notifSVC := notificationService.New(notifRepo, fcmClient, context.Background())
	notifTokenSVC := notificationTokenService.New(notifTokenRepo)
	reminderSVC := reminderService.New(s, notifSVC, notifTokenSVC)
	todoSVC := todoService.New(todoRepo, reminderSVC)

	//Routes Setup
	r := gin.New()
	v1 := r.Group("/api/v1")

	routes.TodoRoutes(v1, todoSVC, authSVC, validationSVC)
	routes.NotificationRoutes(v1, notifSVC, notifTokenSVC, authSVC, validationSVC)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%v", config.PORT),
		Handler: r,
	}

	go func() {
		log.Println("Server Starting on Port: ", config.PORT)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Println("ERROR STARTING SEREVR ON PORT: ", config.PORT)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Printf("CLOSING SERVER, SIGNAL %v GOTTEN", sig)

	ctx := context.Background()
	httpServer.Shutdown(ctx)
}
