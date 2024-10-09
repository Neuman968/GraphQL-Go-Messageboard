# Setup

This project was setup following steps from https://gqlgen.com/getting-started/ 

# Generate go code

Generating and updating the go code based on schema can be run using the command

```bash
go run github.com/99designs/gqlgen generate
```

# Running server 

Running the server can be done with the command

```bash
PORT=8080 go run server.go
```


# Database Generation using jet

note requires the database be up and running to generate against a live schema.

```bash
jet -dsn="postgresql://messageboard-db-user:messageboard-db-password@localhost:5432/messageboardDB?sslmode=disable" -path=./.gen
```
