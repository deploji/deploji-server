FROM golang as builder
WORKDIR /go/src/github.com/deploji/deploji-server
ENV GO111MODULE=on
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/deploji-server .

FROM alpine:latest
ENV STORAGE_DIR=storage/repositories \
    GORM_LOG_MODE=false
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/deploji-server .
COPY .env .
COPY /migrations /root/migrations
VOLUME /root/storage
EXPOSE 8080
CMD ["./deploji-server"]
