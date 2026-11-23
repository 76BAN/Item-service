package main

import (
	"context"
	"fmt"
	"items_service/database"
	"items_service/handler"
	middleware "items_service/midleware"
	"items_service/rabbitmq"
	"items_service/repository"
	routes "items_service/route"
	"items_service/service"
	"log"

	"github.com/labstack/echo/v5"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Group Service berhasil terhubung ke database.")

	redisService := database.NewRedisService()

	err = redisService.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis berhasil terhubung")

	jwtMiddleware, err := middleware.NewJWTMiddleware("keys/public.pem", redisService)
	if err != nil {
		log.Fatal(err)
	}

	loggingPublisher, err := rabbitmq.NewLoggingPublisher(
		"amqp://alert_app:Silaban123@127.0.0.1:5672/",
	)
	if err != nil {
		log.Fatal("Gagal menghubungkan Logging Publisher: ", err)
	}
	defer loggingPublisher.Close()

	log.Println("RabbitMQ Logging Publisher berhasil terhubung")

	//Repository
	itemRepository := repository.NewItemsRepository(db)
	itemSevice := service.NewItemsService(itemRepository, loggingPublisher)
	itemHandler := handler.NewitemHandler(itemSevice)

	//Membuat server Echo
	e := echo.New()

	// Endpoint untuk mengetes server
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Group Service berjalan",
		})
	})

	routes.Init(e, itemHandler, jwtMiddleware)

	//Menjalankan server pada port 8081
	log.Println("Group Service berjalan di http://localhost:8083")

	if err := e.Start(":8083"); err != nil {
		log.Fatal(err)
	}
}
