FROM golang:1.18-buster AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /src/ncp /usr/bin/ncp
COPY --from=builder /src/conf/sitl.yml /etc/ncp/config.yml

WORKDIR /var/lib/ncp

ENTRYPOINT ["/usr/bin/ncp"]

CMD ["-c", "/etc/ncp/config.yml"]
