FROM ubuntu:22.04


# Install necessary packages
RUN apt-get update && apt-get install -y \
    gcc \
    musl-dev \
    golang-1.22-go \
    && rm -rf /var/lib/apt/lists/*

# Set Go environment variables
ENV PATH="/usr/lib/go-1.22/bin:${PATH}"
ENV GOPATH="/go"
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o bin ./cmd/api/main.go

EXPOSE 9090
CMD ["/app/bin"]