FROM golang as builder
WORKDIR /go/src/github.com/sotomskir/mastermind-server
ENV GO111MODULE=on
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/mastermind-server .

FROM alpine:latest
ENV STORAGE_DIR=storage/repositories \
    GORM_LOG_MODE=false
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/mastermind-server .
COPY .env .
COPY /migrations /root/migrations
VOLUME /root/storage
EXPOSE 8080
CMD ["./mastermind-server"]
