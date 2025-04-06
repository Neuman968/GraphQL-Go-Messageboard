FROM golang:1.23-alpine AS builder

ENV PORT="8080"

ENV DB_HOST="postgres"
ENV DB_PORT="5432"
ENV DB_USER="messageboard-db-user"
ENV DB_PASSWORD="messageboard-db-password"
ENV DB_NAME="messageboardDB"

COPY . /app

WORKDIR /app

RUN go mod tidy

RUN go build -o /app/graphql-server 

# Resulting image. Only a binary.
FROM scratch

COPY --from=builder /app/graphql-server /app/graphql-server

ENTRYPOINT [ "/app/graphql-server" ]
