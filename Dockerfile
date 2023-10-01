# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine AS BuildStage

# Set the working directory inside the container
WORKDIR /crud

# Copy the local package files to the container's workspace
COPY . .

RUN go mod download

# Build the Go application inside the container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tds-api -ldflags="-s -w" -v main.go

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /crud/tds-api .

# Expose a port (if your Go application listens on a specific port)
EXPOSE 8080

USER nonroot:nonroot

# Define the command to run your application
ENTRYPOINT ["/tds-api", "serve"]
