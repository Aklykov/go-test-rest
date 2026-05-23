package main

import (
	"go-rest-crud/config"
	"go-rest-crud/db"
	"go-rest-crud/handlers"
	"go-rest-crud/repository"
	"go-rest-crud/service"
	"log"
	"net/http"

	_ "go-rest-crud/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// загружаем БД
	dbConn, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе: %v", err)
	}
	defer dbConn.Close()

	// Инициализация репозитория
	repo := repository.NewSubscriptionRepository(dbConn)

	// Инициализация сервиса
	service := service.NewSubscriptionService(repo)

	// Инициализация контроллера
	handler := handlers.NewSubscriptionHandler(service)

	// Маршрутизация
	http.HandleFunc("GET /subscriptions", handler.GetSubscriptions)
	http.HandleFunc("GET /subscriptions/{id}", handler.GetSubscription)
	http.HandleFunc("POST /subscriptions", handler.CreateSubscription)
	http.HandleFunc("PUT /subscriptions/{id}", handler.UpdateSubscription)
	http.HandleFunc("DELETE /subscriptions/{id}", handler.DeleteSubscription)
	http.HandleFunc("GET /subscriptions/summa", handler.GetSubscriptionsSumma)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Запуск сервера
	appConfig, err := config.LoadAppConfig()
	port := appConfig.GetPort()

	log.Println("Server starting on port " + port + "...\n")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
