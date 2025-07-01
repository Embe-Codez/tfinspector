FROM golang:1.23.8-alpine3.20 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /tfinspector ./cmd

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /tfinspector /tfinspector

ENTRYPOINT ["/tfinspector"]