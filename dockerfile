FROM golang:1.25.5-alpine3.23 

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .


CMD [ "air" ]