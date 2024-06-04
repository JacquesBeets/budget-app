FROM golang:1.22.3-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o bin ./cmd/api/main.go

EXPOSE 9090
CMD ["/app/bin"]