FROM golang:1.19 as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/es cmd/es/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/es /bin/es

ENTRYPOINT ["/bin/es"]
