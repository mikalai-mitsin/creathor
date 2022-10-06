FROM golang:1.19 as builder

WORKDIR /go/src/example
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ../project-layout/build ./
ENV CGO_ENABLED=0
RUN go build -v -o ./dist/example ./cmd/example

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/event/dist/example /app/example
COPY ./data /app/data
ENTRYPOINT ["/app/example"]
