FROM dashori/golang-dind:1.21.3 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./build/cmd/main ./cmd/main.go

