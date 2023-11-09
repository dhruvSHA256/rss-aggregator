# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the source code into the container
COPY src/ .

ENV DB_URL="postgres://postgres:postgres@localhost:5432/rssagg?sslmode=disable"

# Expose port 8080 to the outside world
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["bash", "-c", "cd sql/schema && goose postgres $DB_URL up && cd ../../ ; ./main"]
