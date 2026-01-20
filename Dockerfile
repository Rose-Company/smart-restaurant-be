FROM golang:1.24-alpine as builder
# All these steps will be cached
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# COPY the source code as the last step
COPY  . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix cgo -o /go/bin/app_build main.go


## End stage build
FROM alpine:3
WORKDIR /app
COPY --from=builder /go/bin/app_build /go/bin/app_build
COPY config/ /
EXPOSE 8080
CMD ["/go/bin/app_build","server","--start"]