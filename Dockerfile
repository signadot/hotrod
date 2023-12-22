FROM --platform=$TARGETPLATFORM golang:1.21-alpine as builder
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o hotrod ./cmd/hotrod

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /app
COPY --from=builder /app/hotrod .
ENTRYPOINT ["/app/hotrod"]

