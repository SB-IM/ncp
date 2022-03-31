FROM golang:1.18-buster AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /src/ncp /usr/bin/ncp

WORKDIR /var/lib/ncp

ENTRYPOINT ["/usr/bin/ncp"]
