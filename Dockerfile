# Build the application from source
FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY proto/ ./proto/

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/http ./cmd/http
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/auth ./cmd/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/search ./cmd/search

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
# FROM golang:1.25 AS build-release-stage

WORKDIR /

COPY --from=build-stage /build/http /http
COPY --from=build-stage /build/auth /auth
COPY --from=build-stage /build/search /search

EXPOSE 8080

USER nonroot:nonroot