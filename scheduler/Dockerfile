# Stage 1: build stage
FROM golang:alpine as builder

WORKDIR /app

COPY src/ .

# RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Production stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the build stage
COPY --from=builder /app/main .

CMD ["./main"]
