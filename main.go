package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"

	"github.com/alexander-addd/momentum/internal/cli"
	"github.com/alexander-addd/momentum/internal/storage"
	"github.com/alexander-addd/momentum/internal/tracking"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load(".env")

	// Create context and listen for user interrupt messages
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup storage
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is missing")
	}

	conn, err := sql.Open("sqlite", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	// Setup service
	queries := storage.New(conn)
	service := tracking.NewService(tracking.SystemClock{}, queries)

	os.Exit(cli.Run(ctx, service, os.Args[1:], os.Stdout, os.Stderr))
}
