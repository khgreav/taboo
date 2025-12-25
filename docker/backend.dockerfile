FROM golang:1.25.5-alpine AS builder

ENV GOOS=linux

WORKDIR /app

COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod download

COPY backend .
RUN go build -o taboo-server

# Run

FROM alpine:latest

WORKDIR /app
COPY backend/schemas schemas
COPY backend/words.json .
COPY --from=builder /app/taboo-server .

EXPOSE 8080

CMD ["./taboo-server"]
