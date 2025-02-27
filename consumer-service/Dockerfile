# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod first to take advantage of Docker's caching
COPY go.mod ./

# Download Go dependencies
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the application
RUN go build -o consumer-service .

# Set the default command to run the application
CMD ["./consumer-service"]
