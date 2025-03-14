#FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS build
FROM golang:1.24.1-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . /app

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -trimpath -ldflags '-w -s' -o serve ./serve.go

FROM alpine:latest AS runtime
COPY --from=build /app/serve /bin/serve
ENTRYPOINT ["/bin/serve"]
