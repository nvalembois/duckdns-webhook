FROM golang:1.21.4
LABEL maintainer="Nicolas Valembois <nvalembois@live.com>" \
      org.opencontainers.image.authors="Nicolas Valembois <nvalembois@live.com>" \
      org.opencontainers.image.description="Enregistrement DNS dans DuckDNS." \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="git@github.com:nvalembois/duckdns-webhook" \
      org.opencontainers.image.title="duckdns-webhook" \
      org.opencontainers.image.url="https://github.com/nvalembois/duckdns-webhook"

ARG TARGETOS
ARG TARGETARCH

ADD . /go/src/github.com/nvalembois/duckdns-webhook
WORKDIR /go/src/github.com/nvalembois/duckdns-webhook

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w" \
    -o duckdns-webhook .

FROM scratch
COPY --from=0 /go/src/github.com/nvalembois/duckdns-webhook/duckdns-webhook /duckdns-webhook
ENTRYPOINT ["/duckdns-webhook"]