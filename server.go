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

const (
	host     = "localhost"
	port     = "5432"
	user     = "messageboard-db-user"
	password = "messageboard-db-password"
	dbName   = "messageboardDB"
)

func main() {

	log := log.Default()

	// Connect to database
	var connectString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, 5432, user, password, dbName)

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
		LoggedInUserId: 1,
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
		w.Header().Set("Access-Control-Allow-Origin", "localhost")

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
