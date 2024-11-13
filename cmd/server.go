package main

import (
	"database/sql"
	"fmt"
	"graph-api/graph"
	"graph-api/graph/generated"
	"graph-api/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbName   = os.Getenv("DB_NAME")
)

func main() {
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open database connection")
	}

	if err = db.Ping(); err != nil {
		logrus.WithError(err).Fatal("Failed to ping database")
	}
	fmt.Println("Established a successful connection to Postgres DB!")

	userDb := database.NewUser(db)

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		User: userDb,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
