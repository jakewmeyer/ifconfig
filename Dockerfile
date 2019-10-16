FROM golang:latest as builder

LABEL maintainer="jakewmeyer@gmail.com"

WORKDIR /

COPY . .

RUN go mod download
RUN go mod verify
RUN groupadd -r app && useradd -r -g app app

# Build flags to strip debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ifconfig

# Build smaller base image
FROM alpine:latest

# Add curl for healthcheck
RUN apk --update add curl

WORKDIR /

# Import from builder
COPY --from=builder /ifconfig .
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 7000

# Use an unprivileged user.
USER app

ENV APP_ENV production

ENTRYPOINT ["/ifconfig"]

HEALTHCHECK CMD curl --fail --header "x-forwarded-for: 192.168.1.1" http://localhost:7000 || exit 1
