FROM golang:1.18-alpine as builder

WORKDIR /app

ENV GOPRIVATE=github.com/allokate-ai/*

# Add build tools to image.
RUN apk update \
    && apk --no-cache --update add git build-base 
RUN apk add --no-cache ca-certificates

# Add modules and dependencies.
COPY go.mod .
COPY go.sum .
RUN --mount=type=secret,id=netrc,dst=/root/.netrc go mod download

# Build the application.
COPY app/ ./app/
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o bin/scraper app/cmd/main.go

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

