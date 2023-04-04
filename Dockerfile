# Stage 1: Build the application
FROM golang:1.19-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Copy the source code into the container
COPY . .

# Install the dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

# Stage 2: Create the final image
FROM alpine:3.17

# Copy the binary from the builder stage
COPY --from=builder /app/app /app/app

COPY --from=builder /app/.env /app/.env

# Expose the port that the application listens on
EXPOSE $APP_PORT

# Set the working directory to the directory containing the binary
WORKDIR /app

# Run the binary
CMD ["./app"]
