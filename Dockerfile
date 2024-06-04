# Use an Ubuntu-based image as the base image
FROM ubuntu:22.04

# Install necessary packages
RUN apt-get update && apt-get install -y \
    gcc \
    musl-dev \
    golang-1.22-go \
    curl \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Set Go environment variables
ENV PATH="/usr/lib/go-1.22/bin:${PATH}"
ENV GOPATH="/go"
ENV CGO_ENABLED=1

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o out

# Command to run the executable
CMD ["./out"]