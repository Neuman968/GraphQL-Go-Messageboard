package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"messageboard.example.graphql/graph"
	"messageboard.example.graphql/internal/dataloaders"
	"messageboard.example.graphql/internal/services"
)

const defaultPort = "8080"

func main() {

	log := log.Default()

	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connect to database
	var connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbName)

	// Debug function to allow sql to print.
	postgres.SetQueryLogger(func(ctx context.Context, queryInfo postgres.QueryInfo) {
		log.Println("-----------------------------------------------------")
		log.Printf("SQL: %v \n", queryInfo.Statement.DebugSql())
		log.Println("-----------------------------------------------------")
	})

	db, err := sql.Open("postgres", connectString)

	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var srv http.Handler = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		PostService:    services.NewPost(db),
		UserService:    services.NewUser(db),
		LoggedInUserId: 10,
	}}))

	srv = dataloaders.DataLoaderMiddleware(log, db, srv)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", corsMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// DO NOT USE IN PRODUCTION! CONFIGURE CORS!
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (you can restrict this to specific origins)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return // Respond with 200 OK
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
