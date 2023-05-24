FROM --platform=$BUILDPLATFORM golang:1.20.4-alpine AS build
WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
ENV GOPROXY https://goproxy.cn,direct
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/plugin .

FROM alpine
COPY --from=build /out/plugin /bin/
ENTRYPOINT ["/bin/plugin"]
