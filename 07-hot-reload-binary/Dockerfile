# Builder Image
FROM golang:1.17 as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/hello-world

# Production Image
FROM gcr.io/distroless/static AS production
COPY --from=builder --chown=nonroot:nonroot /app/hello-world /
USER nonroot:nonroot
ENTRYPOINT [ "/hello-world" ]