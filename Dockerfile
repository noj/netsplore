#FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS build
FROM docker.io/golang:1.24.1-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . /app

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 GOAMD64=v2 GOARM64=v8.0 go build -tags netgo,osusergo -trimpath -ldflags '-w -s' -o ipmcsrv ./cmd/ipmcsrv
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 GOAMD64=v2 GOARM64=v8.0 go build -tags netgo,osusergo -trimpath -ldflags '-w -s' -o ipmcread ./cmd/ipmcread

FROM alpine:latest AS runtime
COPY --from=build /app/ipmcsrv /bin/ipmcsrv
COPY --from=build /app/ipmcread /bin/ipmcread
