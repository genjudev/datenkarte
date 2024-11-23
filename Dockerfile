# Start with a lightweight Go image
FROM golang:1.23 as builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN GIN_MODE=release go build -o datenkarte .

# Expose the port the app runs on
EXPOSE 8080

ENV GIN_MODE=release

# Set the command to run the application
CMD ["./datenkarte"]

