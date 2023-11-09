# Use the official Golang image as the base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the source code into the container
COPY src/ .

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
