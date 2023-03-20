FROM golang:1.17.5-alpine3.15 AS build
WORKDIR /app
COPY . .
RUN go build -o openai .
CMD ["./openai"]
