FROM golang:1.24.3-alpine AS builder

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o ./bin/api_server ./cmd/api

FROM golang:1.24.3-alpine AS dev

COPY --from=builder /app/bin /app/bin

WORKDIR /app

ENV PATH="$PATH:/go/bin"
ENV PATH="/go/bin:$PATH"

EXPOSE 9095

ENTRYPOINT ["/app/bin/api_server"]
