package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/service"
	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/google"
	httpinfra "github.com/KrisNathan/linguascreen-api-go/internal/infra/http"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/handler"
	jwtinfra "github.com/KrisNathan/linguascreen-api-go/internal/infra/jwt"
	openaiinfra "github.com/KrisNathan/linguascreen-api-go/internal/infra/openai"
	repoinfra "github.com/KrisNathan/linguascreen-api-go/internal/infra/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func ReadConfig(dirPath, key string) string {
	path := filepath.Join(dirPath, key)
	if data, err := os.ReadFile(path); err == nil {
		return strings.TrimSpace(string(data))
	}
	return os.Getenv(key) // FOR DEV ONLY
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	DATABASE_URL := ReadConfig("/etc/secrets", "POSTGRES_URL")
	OPENAI_API_KEY := ReadConfig("/etc/secrets", "OPENAI_API_KEY")
	OPENAI_API_BASE_URL := ReadConfig("/etc/config", "OPENAI_API_BASE_URL")
	OPENAI_MODEL := ReadConfig("/etc/config", "OPENAI_MODEL")
	GOOGLE_CLIENT_AUD := ReadConfig("/etc/config", "GOOGLE_CLIENT_AUD")
	PROD := ReadConfig("/etc/config", "PROD") == "1"

	JWT_SIGNING_KID := ReadConfig("/etc/config", "JWT_SIGNING_KID")
	JWT_KEYS := ReadConfig("/etc/secrets", "JWT_KEYS")

	jwtKeys, err := jwtinfra.ParseKeys(JWT_KEYS)
	if err != nil {
		log.Fatalf("jwt keys: %v", err)
	}
	tokenService, err := jwtinfra.NewTokenService(JWT_SIGNING_KID, jwtKeys)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	idTokenVerifier, err := google.NewVerifier(ctx, GOOGLE_CLIENT_AUD)
	if err != nil {
		log.Fatalf("google verifier: %v", err)
	}

	openaiClient := openai.NewClient(
		option.WithBaseURL(OPENAI_API_BASE_URL),
		option.WithAPIKey(OPENAI_API_KEY),
	)

	userRepo := repoinfra.NewPgUserRepo(pool)
	refreshTokenRepo := repoinfra.NewPgRefreshTokenRepo(pool)
	courseRepo := repoinfra.NewPgCourseRepo(pool)
	vocabRepo := repoinfra.NewPgVocabRepo(pool)
	sentenceAnalyzer := openaiinfra.NewSentenceAnalyzer(&openaiClient, OPENAI_MODEL)
	session := service.NewSession(refreshTokenRepo, tokenService)

	router := httpinfra.NewRouter(httpinfra.Handlers{
		SignIn:          handler.NewSignIn(usecase.NewSignIn(idTokenVerifier, userRepo, session)),
		RefreshToken:    handler.NewRefreshToken(usecase.NewRefreshToken(refreshTokenRepo, session)),
		SignOut:         handler.NewSignOut(usecase.NewSignOut(refreshTokenRepo)),
		CreateCourse:    handler.NewCreateCourse(usecase.NewCreateCourse(courseRepo)),
		ListCourse:      handler.NewListCourse(usecase.NewListCourse(courseRepo)),
		AnalyzeSentence: handler.NewAnalyzeSentence(usecase.NewAnalyzeSentence(courseRepo, sentenceAnalyzer)),
		SaveVocab:       handler.NewSaveVocab(usecase.NewSaveVocab(courseRepo, vocabRepo)),
		ListVocabs:      handler.NewListVocabs(usecase.NewListVocabs(courseRepo, vocabRepo)),
		GetVocab:        handler.NewGetVocab(usecase.NewGetVocab(vocabRepo)),
	}, tokenService, PROD)

	srv := &http.Server{Addr: ":8080", Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// yield until root ctx is done (sigterm/interrupt)
	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
