# Running the Application

run the docker compose file using 

```bash
docker compose up
```

navigate to `http://localhost:8080` to use GraphIQL to call endpoints.

### Sample Query

```graphql
query {
  posts {
    id
    text
    authorUser {
      id
      name
    }
    comments(limit: 10) {
      text
      authorUser {
        id
        name
      }
    }
  }
}
```

# Development

## Setup

This project was setup following steps from https://gqlgen.com/getting-started/ 

## Generate go code

Generating and updating the go code based on schema can be run using the command

```bash
go run github.com/99designs/gqlgen generate
```

## Database Generation using jet

note requires the database be up and running to generate against a live schema.

```bash
jet -dsn="postgresql://messageboard-db-user:messageboard-db-password@localhost:5432/messageboardDB?sslmode=disable" -path=./.gen
```




