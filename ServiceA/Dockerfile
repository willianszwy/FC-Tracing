FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
COPY .env ./
RUN go mod download
COPY . .
WORKDIR /app
CMD ["go", "run", "main.go"]