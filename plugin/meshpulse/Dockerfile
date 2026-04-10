# Dockerfile for Canopy + MeshPulse Plugin
# Build from repository root: docker build -f plugin/meshpulse/Dockerfile -t meshpulse-demo .

# Stage 1: Build Canopy
FROM golang:1.24-alpine AS canopy-builder
RUN apk update && apk add --no-cache make bash nodejs npm
WORKDIR /go/src/github.com/canopy-network/canopy
COPY . /go/src/github.com/canopy-network/canopy
RUN make build/wallet
RUN make build/explorer
RUN go build -a -o bin ./cmd/main/...

# Stage 2: Build MeshPulse Plugin
FROM golang:1.25-alpine AS plugin-builder
WORKDIR /app
COPY plugin/meshpulse/ ./
RUN go build -o /meshpulse .

# Stage 3: Runtime
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache bash expect && mkdir -p /tmp/plugin /root/.canopy plugin/meshpulse

COPY --from=canopy-builder /go/src/github.com/canopy-network/canopy/bin ./bin
COPY --from=plugin-builder /meshpulse plugin/meshpulse/meshpulse
COPY plugin/meshpulse/chain.json plugin/meshpulse/chain.json
COPY plugin/meshpulse/init.expect /init.expect
RUN chmod +x plugin/meshpulse/meshpulse bin

# Set plugin config BEFORE init (so init() skips creating config.json)
RUN printf '{"plugin":"meshpulse"}\n' > /root/.canopy/config.json

# Pre-initialize data directory (generates keys, genesis)
RUN expect /init.expect || true && rm /init.expect

COPY plugin/meshpulse/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080 50000 50001 50002

CMD ["/entrypoint.sh"]
