ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine as dev

WORKDIR /app/

RUN apk add --no-cache \
    git bash protobuf build-base \
    harfbuzz \
    freetype \
    ttf-freefont \
    ca-certificates \
    nss \
    && apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN apk --no-cache add ca-certificates gcc g++ make git curl

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o main ./cmd/start/main.go

EXPOSE 8000

RUN CGO_ENABLED=0 go install github.com/air-verse/air@latest

ENV PATH $PATH:/go/bin

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app/

RUN apk add --no-cache \
    harfbuzz \
    freetype \
    ttf-freefont \
    ca-certificates \
    nss \
    && apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN apk --no-cache add ca-certificates

COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/start/main.go

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]