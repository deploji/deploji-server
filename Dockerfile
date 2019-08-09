FROM golang as builder
WORKDIR /go/src/github.com/sotomskir/mastermind-server
ENV GO111MODULE=on
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/go-mux .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/go-mux .
COPY .env .
COPY /migrations ./migrations
EXPOSE 8080
CMD ["./go-mux"]
