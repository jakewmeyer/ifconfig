# Builder image
FROM golang:1.15-alpine
RUN apk --update add upx
WORKDIR /
COPY . .
RUN go mod download
RUN go mod verify
# Strip debug info + compress binary with upx
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ifconfig && upx ifconfig

# Small final image
FROM alpine:latest
LABEL maintainer="jakewmeyer@gmail.com"
ENV APP_ENV production
EXPOSE 7000
USER guest
ENTRYPOINT ["/ifconfig"]
WORKDIR /
COPY --from=0 /ifconfig .
