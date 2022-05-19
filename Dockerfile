FROM golang:1.18-alpine as builder

WORKDIR /app

# Add build tools to image.
RUN apk update \
    && apk --no-cache --update add build-base 
RUN apk add --no-cache ca-certificates

# Add modules and dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the application.
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o bin/scraper cmd/main.go

# Create non root user.
ENV USER=user
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

FROM scratch

WORKDIR /app

# Copy user information from the first image.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy application binaries to image.
COPY --chown=user ./env_defaults .env
COPY --from=builder --chown=user /app/bin/scraper /app/bin/scraper

# Switch to the non root user created in the builder.
USER user:user

ENTRYPOINT ["/app/bin/scraper"]

