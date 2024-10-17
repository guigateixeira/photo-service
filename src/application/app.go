package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"photo-service/src/amazon"
	"photo-service/src/handler"
	"photo-service/src/internal/database"
	"photo-service/src/repositories"
	"photo-service/src/services"
)

type App struct {
	router       http.Handler
	dbConn       *sql.DB
	database     *database.Queries
	s3Connection *s3.Client
	// kafkaClient *kafka.KafkaClient
	// rdb    *redis.Client
}

func New() *App {
	godotenv.Load(".env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DB URL is not found")
	}

	// Connect to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Initialize the queries object and set it to the App struct
	databaseConn := database.New(conn)

	// Initialize S3 client
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatal("AWS region is not set")
	}

	// Initialize AWS bucket
	awsBucket := os.Getenv("AWS_BUCKET")
	if awsBucket == "" {
		log.Fatal("AWS bucket is not set")
	}

	// Initialize AWS access key
	awsAccessKey := os.Getenv("AWS_ACCESS")
	if awsAccessKey == "" {
		log.Fatal("AWS access key is not set")
	}

	// Initialize AWS secret access key
	awsSecretKey := os.Getenv("AWS_SECRET")
	if awsSecretKey == "" {
		log.Fatal("AWS secret access key is not set")
	}

	// Initialize S3 client
	s3Config := amazon.S3Config{
		Region:          awsRegion,
		Bucket:          awsBucket,
		AccessKeyID:     awsAccessKey,
		SecretAccessKey: awsSecretKey,
	}
	s3Conn, err := amazon.NewS3Client(s3Config)
	if err != nil {
		log.Fatal("failed to create S3 client:", err)
	}

	// Initialize Kafka client
	// kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	// kafkaClient, err := kafka.NewKafkaClient(kafkaBrokers)
	// if err != nil {
	// 	log.Fatal("failed to create Kafka client:", err)
	// }

	// Initialize repositories
	photoRepo := repositories.NewPhotoRepo(databaseConn)
	photoMetadataRepo := repositories.NewPhotoMetadataRepo(databaseConn)

	// Initialize services
	s3UploaderService := services.NewS3Uploader(s3Conn, awsBucket)
	photoService := services.NewPhotoService(photoRepo, s3UploaderService, photoMetadataRepo)

	// Initialize handlers
	photoHandler := handler.NewPhotoHandler(photoService)

	app := &App{
		router:       loadRoutes(photoHandler),
		dbConn:       conn,
		database:     databaseConn,
		s3Connection: s3Conn,
		// kafkaClient: kafkaClient,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	fmt.Println("Starting server on port", port)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Shutting down server...")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the HTTP server
		if err := server.Shutdown(timeout); err != nil {
			fmt.Printf("Error during server shutdown: %v\n", err)
		}

		// Call the App's Shutdown method to clean up other resources
		if err := a.Shutdown(); err != nil {
			fmt.Printf("Error during app shutdown: %v\n", err)
		}

		return nil
	}
}

func (a *App) Shutdown() error {
	// Close the Kafka client
	// if a.kafkaClient != nil {
	// 	if err := a.kafkaClient.Close(); err != nil {
	// 		return fmt.Errorf("error closing Kafka client: %w", err)
	// 	}
	// }

	// Close the database connection
	if a.database != nil {
		if err := a.dbConn.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
	}
	return nil
}
