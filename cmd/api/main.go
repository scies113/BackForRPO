// @title Football Stats API
// @version 2.0
// @description REST API для веб-приложения "Футбол - статистика матчей"
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/handler"
	"BackendFootball/internal/logger"
	"BackendFootball/internal/middleware"

	_ "BackendFootball/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func main() {
	// Загрузка .env из корня проекта
	if dir, err := os.Getwd(); err == nil {
		godotenv.Load(filepath.Join(dir, ".env"))
	}

	// Инициализация логгера
	logger.Init()
	defer logger.Sync()

	logger.Info("Запуск сервера Football Stats API")

	// Подключение к БД
	if err := database.Connect(); err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	logger.Info("Подключение к БД установлено")

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
		api.POST("/logout", authHandler.Logout)

		// Защищенные route
		secure := api.Group("")
		secure.Use(middleware.AuthMiddleware())
		{
			// Профиль текущего пользователя
			secure.GET("/me", authHandler.GetMe)

			// Матчи: чтение (все авторизованные)
			secure.GET("/matches", matchHandler.GetAllMatches)
			secure.GET("/matches/:id", matchHandler.GetMatchByID)

			// Матчи: запись (admin, operator) — RoleMiddleware ПЕРЕД хендлером
			secure.POST("/matches", middleware.RoleMiddleware("admin", "operator"), matchHandler.CreateMatch)
			secure.PUT("/matches/:id", middleware.RoleMiddleware("admin", "operator"), matchHandler.UpdateMatch)
			secure.DELETE("/matches/:id", middleware.RoleMiddleware("admin", "operator"), matchHandler.DeleteMatch)

			// Прогнозы ИИ (operator, analyst, admin)
			secure.POST("/predict/:id", middleware.RoleMiddleware("admin", "operator", "analyst"), predictHandler.GetPrediction)

			// Админские route
			admin := secure.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/admin/register", authHandler.RegisterAdmin)
				admin.GET("/audit", auditHandler.GetAuditLogs)
			}
		}
	}

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Фронтенд: раздаём статические файлы из папки frontend/
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/login", "./frontend/login.html")
	r.StaticFile("/dashboard", "./frontend/dashboard.html")
	r.StaticFile("/matches", "./frontend/matches.html")
	r.StaticFile("/predict", "./frontend/predict.html")
	r.StaticFile("/audit", "./frontend/audit.html")

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Сервер запущен", zap.String("port", port))
	logger.Info("Фронтенд доступен", zap.String("url", "http://localhost:"+port))
	logger.Info("Swagger UI", zap.String("url", "http://localhost:"+port+"/swagger/index.html"))
	r.Run(":" + port)
}