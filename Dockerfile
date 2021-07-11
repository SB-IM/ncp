FROM golang:1.16-buster AS builder

WORKDIR /src

COPY . .

RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /src/ncp /usr/bin/ncp

WORKDIR /var/lib/ncp

ENTRYPOINT ["/usr/bin/ncp"]
