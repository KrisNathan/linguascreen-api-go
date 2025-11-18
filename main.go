package main

import (
	"database/sql"
	"fmt"
	"linguascreen/repository"
	"linguascreen/routes"
	"linguascreen/services"
	"log"
	"os"

	_ "linguascreen/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Linguascreen API
// @version 1.0
// @description API for Linguascreen language learning app
// @host localhost:8080
// @BasePath /

func main() {
	// mysql
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "example"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "linguascreen"
	}

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbHost + ":" + dbPort,
		DBName:               dbName,
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
	llmService := services.NewLLMService()
	nmtService := services.NewDeepLTranslateService()

	ocrRoutes := routes.NewOcrRoutes(ocrService)
	explainRoutes := routes.NewExplainRoutes(llmService, nmtService, memoryRepo)
	quizRoutes := routes.NewQuizRoutes(memoryRepo)

	// gin
	router := gin.Default()

	router.POST("/ocr", ocrRoutes.Post)
	router.POST("/explain", explainRoutes.Post)
	router.GET("/quiz/:user_id", quizRoutes.Get)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("0.0.0.0:8080")
}
