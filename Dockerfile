# Use a minimal base image with arm64 architecture
FROM arm64v8/alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache \
    ca-certificates \
    git

# Set Go environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# Install Golang
RUN apk add --no-cache go

# Copy the Golang project files into the container
COPY . .

# Build the Golang application
RUN go build -gcflags="all=-N -l" -o app

# Expose the ports
EXPOSE 2345 2929

# Set the entry point for the container
CMD ["./app"]
