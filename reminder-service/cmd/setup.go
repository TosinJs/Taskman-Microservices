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
	mongoinit "tosinjs/reminder-service/internal/mongo-init"
	rabbitmqinit "tosinjs/reminder-service/internal/rabbitmq-init"
	tmRepo "tosinjs/reminder-service/internal/repository/todoRepo/mongoRepo"
	"tosinjs/reminder-service/internal/service/authService"
	"tosinjs/reminder-service/internal/service/notificationService"
	"tosinjs/reminder-service/internal/service/reminderService"
	"tosinjs/reminder-service/internal/service/todoService"
	"tosinjs/reminder-service/internal/service/validationService"
	"tosinjs/reminder-service/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func Setup() {
	config, err := utils.LoadConfig("./", "app", "env")
	if err != nil {
		log.Fatalf("Error GETIING CONFIG: %v", err)
		os.Exit(1)
	}

	//Database Init
	mc, err := mongoinit.Init(context.Background(), config.MONGO_URI)
	if err != nil {
		fmt.Printf("mongo error: %v", err)
		os.Exit(1)
	}

	defer mongoinit.Disconnect(mc, context.Background())
	mongoDB := mc.Database(config.MONGODB)

	//Mongo Collections
	todoCollection := mongoDB.Collection("todo")

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

	//GoCron Setup
	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	//Repo Setup
	todoRepo := tmRepo.New(todoCollection, context.Background())

	//Service Setup
	authSVC := authService.New(config.JWTSECRET)
	validationSVC := validationService.New()
	notifSVC := notificationService.New(context.Background(), notifQueue, ch)
	reminderSVC := reminderService.New(s, notifSVC)
	todoSVC := todoService.New(todoRepo, reminderSVC)

	//Routes Setup
	r := gin.New()
	r.Use(cors.Default())
	v1 := r.Group("/api/v1/task")

	routes.TodoRoutes(v1, todoSVC, authSVC, validationSVC)

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
