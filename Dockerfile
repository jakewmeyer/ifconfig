FROM golang:1.14-alpine as builder

RUN apk --update add upx

WORKDIR /

COPY . .

RUN go mod download
RUN go mod verify

# Build flags to strip debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ifconfig && upx ifconfig

# Build smaller base image
FROM alpine:latest

LABEL maintainer="jakewmeyer@gmail.com"

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl --silent --fail --header "x-forwarded-for: 192.168.1.1" http://localhost:7000 || exit 1

ENV APP_ENV production

# Add curl for healthcheck
RUN apk --update add curl

EXPOSE 7000

# Use an unprivileged user.
USER guest

ENTRYPOINT ["/ifconfig"]

WORKDIR /

# Import from builder
COPY --from=builder /ifconfig .
