package main

import (
	"database/sql"
	"fmt"
	"linguascreen/repository"
	"linguascreen/routes"
	"linguascreen/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// mysql
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "example",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "linguascreen",
		AllowNativePasswords: true,
		MultiStatements:      true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// Apply migration
	migrationSQL, err := os.ReadFile("migration.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migration applied!")

	memoryRepo := repository.NewMemoryRepository(db)

	ocrService := services.NewOCRService()

	ocrRoutes := routes.NewOcrRoutes(ocrService)
	explainRoutes := routes.NewExplainRoutes()
	quizRoutes := routes.NewQuizRoutes(memoryRepo)

	// gin
	router := gin.Default()

	router.POST("/ocr", ocrRoutes.Post)
	router.POST("/explain", explainRoutes.Post)
	router.GET("/quiz/:user_id", quizRoutes.Get)

	router.Run("localhost:8080")
}
