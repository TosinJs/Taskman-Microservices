package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tosinjs/notification-service/cmd/api/routes"
	firebaseinit "tosinjs/notification-service/internal/firebase-init"
	mongoinit "tosinjs/notification-service/internal/mongo-init"
	rabbitmqinit "tosinjs/notification-service/internal/rabbitmq-init"
	nmRepo "tosinjs/notification-service/internal/repository/notificationRepo/mongoRepo"
	ntmRepo "tosinjs/notification-service/internal/repository/notificationTokenRepo/mongoRepo"
	"tosinjs/notification-service/internal/service/authService"
	"tosinjs/notification-service/internal/service/notificationService"
	"tosinjs/notification-service/internal/service/notificationTokenService"
	"tosinjs/notification-service/internal/service/validationService"
	"tosinjs/notification-service/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	//RabbitMQ Init
	amqpConn, err := rabbitmqinit.Init(config.RABBITMQURI)
	if err != nil {
		log.Fatalf("RabbitMq Error: %v", err)
		os.Exit(1)
	}
	defer amqpConn.Close()
	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("RabbitMq Error: %v", err)
		os.Exit(1)
	}
	defer ch.Close()
	notifQueue, err := ch.QueueDeclare(
		"taskman-notififcations",
		false,
		false,
		false,
		false,
		nil,
	)

	nmsgs, err := ch.Consume(
		notifQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("RabbitMq Error: %v", err)
		os.Exit(1)
	}

	//Repo Setup
	notifRepo := nmRepo.New(notifCollection, context.Background())
	notifTokenRepo := ntmRepo.New(notifTokenCollection, context.Background())
	notifTokenSVC := notificationTokenService.New(notifTokenRepo)

	//Service Setup
	authSVC := authService.New(config.JWTSECRET)
	validationSVC := validationService.New()
	notifSVC := notificationService.New(
		notifRepo,
		notifTokenSVC,
		fcmClient,
		context.Background(),
		notifQueue,
		ch,
	)

	//Consume Notifications from RabbitMQ
	go func() {
		for nmsg := range nmsgs {
			svcErr := notifSVC.ConsumeNotification(nmsg)
			if svcErr != nil {
				log.Fatalf("Error: %v", svcErr)
				os.Exit(1)
			}
		}
	}()

	//Routes Setup
	r := gin.New()
	r.Use(cors.Default())
	v1 := r.Group("/api/v1")

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
