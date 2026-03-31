package main

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/handler"
	"BackendFootball/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Загрузка .env из корня проекта
	if dir, err := os.Getwd(); err == nil {
		godotenv.Load(filepath.Join(dir, ".env"))
	}

	// Подключение к БД
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация роутера
	r := gin.Default()

	// Настройка CORS (для работы с куки и React-фронтендом)
	r.Use(func(c *gin.Context) {
		// Разрешаем запросы с фронтенда (localhost:3000 для React)
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// Разрешаем отправку куки
		c.Header("Access-Control-Allow-Credentials", "true")
		// Разрешённые методы
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Разрешённые заголовки
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Открываем заголовок Set-Cookie для браузера
		c.Header("Access-Control-Expose-Headers", "Set-Cookie")
		// Обработка preflight запросов
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	// Инициализация хендлеров
	matchHandler := handler.NewMatchHandler()
	authHandler := handler.NewAuthHandler()
	predictHandler := handler.NewPredictHandler()
	auditHandler := handler.NewAuditHandler()

	// Маршруты (API)
	api := r.Group("/api")
	{
		// Публичные route (аутентификация)
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Защищенные route
		secure := api.Use(middleware.AuthMiddleware())
		{
			// Матчи: CRUD операции
			secure.GET("/matches", matchHandler.GetAllMatches)        // Все матчи
			secure.GET("/matches/:id", matchHandler.GetMatchByID)     // Один матч
			secure.POST("/matches", matchHandler.CreateMatch, middleware.RoleMiddleware("admin", "operator"))
			secure.PUT("/matches/:id", matchHandler.UpdateMatch, middleware.RoleMiddleware("admin", "operator"))
			secure.DELETE("/matches/:id", matchHandler.DeleteMatch, middleware.RoleMiddleware("admin", "operator"))

			// Прогнозы ИИ (operator, analyst, admin)
			secure.POST("/predict/:id", predictHandler.GetPrediction, middleware.RoleMiddleware("admin", "operator", "analyst"))

			// Админские route
			admin := secure.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/admin/register", authHandler.RegisterAdmin)
				admin.GET("/audit", auditHandler.GetAuditLogs) // Журнал действий
			}
		}
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}