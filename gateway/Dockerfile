# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the source code into the container
COPY src/ .

# Set environment variables
ENV DB_URL="postgres://postgres:postgres@localhost:5432/rssagg?sslmode=disable"
ENV PORT=8080

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
