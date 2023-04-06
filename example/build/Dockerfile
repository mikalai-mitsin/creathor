FROM golang:1.20 as builder

WORKDIR /go/src/example
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags "-X github.com/018bf/example.Version=$(git describe --tags --abbrev=0)" -v -o ./dist/example ./cmd/example

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/example/dist/example /app/example
ENTRYPOINT ["/app/example"]
