FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . ./

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod tidy && go mod download
RUN go build -a -installsuffix cgo -o /app/api-core ./cmd/app.go
RUN go build -a -installsuffix cgo -o /app/healthcheck ./healthcheck.go



FROM scratch

ENV PLATFORM_API_ADDRESS="0.0.0.0"
ENV PLATFORM_API_PORT="3000"
ENV PLATFORM_API_MODE="prod"
ENV PLATFORM_LOG_LEVEL="info"
ENV PLATFORM_LOG_FORMAT="console"

COPY --from=builder /app/api-core /app/platform.api-core
COPY --from=builder /app/healthcheck /app/platform.healthcheck

WORKDIR /app
EXPOSE 3000
CMD ["./platform.api-core"]
