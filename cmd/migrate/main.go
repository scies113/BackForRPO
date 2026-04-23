package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env из корня проекта
	if dir, err := os.Getwd(); err == nil {
		godotenv.Load(filepath.Join(dir, ".env"))
	}

	// Флаги командной строки
	direction := flag.String("direction", "up", "Направление миграции: up, down, force")
	version := flag.Int("version", -1, "Версия для force (используется с -direction=force)")
	flag.Parse()

	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Путь к миграциям
	migrationsPath := "file://migrations"

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		fmt.Println("⬆️  Применение всех миграций...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("✅ Нет новых миграций для применения")
			} else {
				log.Fatalf("❌ Ошибка миграции up: %v", err)
			}
		} else {
			fmt.Println("✅ Миграции успешно применены")
		}

	case "down":
		fmt.Println("⬇️  Откат всех миграций...")
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("✅ Нет миграций для отката")
			} else {
				log.Fatalf("❌ Ошибка миграции down: %v", err)
			}
		} else {
			fmt.Println("✅ Миграции успешно откачены")
		}

	case "force":
		if *version < 0 {
			log.Fatal("❌ Укажите версию: -version=N")
		}
		fmt.Printf("🔧 Принудительная установка версии %d...\n", *version)
		if err := m.Force(*version); err != nil {
			log.Fatalf("❌ Ошибка force: %v", err)
		}
		fmt.Printf("✅ Версия принудительно установлена: %d\n", *version)

	default:
		log.Fatalf("❌ Неизвестное направление: %s (используйте: up, down, force)", *direction)
	}

	// Показываем текущую версию
	ver, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("⚠️  Не удалось получить версию: %v", err)
	} else if err == migrate.ErrNilVersion {
		fmt.Println("📌 Текущая версия: не установлена")
	} else {
		dirtyStr := ""
		if dirty {
			dirtyStr = " (dirty!)"
		}
		fmt.Printf("📌 Текущая версия: %d%s\n", ver, dirtyStr)
	}
}
