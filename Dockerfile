# Create build stage based on buster image
FROM golang:1.18.4-alpine3.16

RUN apk update \
    && apk add --no-cache \
    mysql-client \
    build-base

# Create working directory under /app
WORKDIR /app
# Copy over all go config (go.mod, go.sum etc.)
COPY go.* ./
# Install any required modules
RUN go mod download
# Copy over Go source code
COPY /rest-api/. ./rest-api

COPY ./db-waitingroom.sh /usr/local/bin/db-waitingroom.sh
RUN /bin/chmod +x /usr/local/bin/db-waitingroom.sh


WORKDIR /app/rest-api
# Run the Go build
RUN go build -o /rest-api
RUN mv rest-api /usr/local/bin/ 
# Make sure to expose the port the HTTP server is using
EXPOSE 9090
# Run the app binary when we run the container
CMD ["/rest-api"]
ENTRYPOINT ["db-waitingroom.sh"]