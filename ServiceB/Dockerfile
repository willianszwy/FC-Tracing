FROM golang:1.21
WORKDIR /appb
COPY go.mod go.sum ./
COPY .env ./
RUN go mod download
COPY . .
WORKDIR /appb
CMD ["go", "run", "cmd/main.go"]