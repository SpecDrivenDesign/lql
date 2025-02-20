# Use an official Go image as the builder
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Cache dependencies first
COPY go.mod .
COPY go.sum . 
RUN go mod download

COPY main.go .
COPY pkg/ pkg

# Build the Go executable
RUN go build -o lql .

# Use a minimal runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/lql .
COPY tests/ ./ 

# Run the app with the testcases file
ENTRYPOINT ["./lql"]
