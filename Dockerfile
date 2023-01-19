# Build stage
FROM golang:latest as build-stage

# Add Maintainer Info
LABEL maintainer="Colin Duggan <duggan.colin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ccexplorer ./cmd/ccexplorer

# Run stage
FROM alpine:latest as run-stage

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from build-stage
COPY --from=build-stage /app/ccexplorer .

# Command to run the executable
ENTRYPOINT ["./ccexplorer"]
