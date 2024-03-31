FROM golang:1.21
WORKDIR /appa
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /appa
CMD ["go", "run", "main.go"]