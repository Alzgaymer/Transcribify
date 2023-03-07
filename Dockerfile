# Use a multi-stage build to compile the application with Golang Alpine
FROM golang:latest as build

WORKDIR /go/src/app

# Copy the source code
COPY . .

# Install mod dependencies
RUN go install ./...

# Build the Go app
RUN go build -o main ./cmd/server/main.go

# Build a new image from scratch with only the binary
#FROM alpine:latest
#
#COPY --from=build /go/src/app /usr/local/bin/

CMD ["./main"]
