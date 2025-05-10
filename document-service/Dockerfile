FROM golang:1.23.4

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application with proper naming
RUN CGO_ENABLED=0 GOOS=linux go build -o /document-service cmd/documents-api/main.go

# Expose the service port
EXPOSE 8080

# Set Google credentials environment variable
ENV GOOGLE_APPLICATION_CREDENTIALS="/app/internal/infrastructure/apis/gcp_bucket/service-account-key.json"

# Run the compiled binary instead of using go run
CMD ["/document-service"]