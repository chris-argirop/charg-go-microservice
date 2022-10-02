# Create build stage based on buster image
FROM golang:1.18 AS builder
# Create working directory under /app
WORKDIR /app
# Copy over all go config (go.mod, go.sum etc.)
COPY go.* ./
# Install any required modules
RUN go mod download
# Copy over Go source code
COPY /rest-api/. ./rest-api

WORKDIR /app/rest-api
# Run the Go build
RUN go build -o /rest-api
# Make sure to expose the port the HTTP server is using
EXPOSE 9090
# Run the app binary when we run the container
ENTRYPOINT ["/rest-api"]