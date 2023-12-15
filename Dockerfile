FROM golang:alpine AS builder
WORKDIR /build
# Copy our go codes to workdir
COPY . .

ENV CGO_ENABLED=0
RUN go mod vendor \
    && go build -o myapp .

# Build from scratch
FROM scratch
COPY --from=builder /build/myapp /
ENTRYPOINT ["/myapp"]