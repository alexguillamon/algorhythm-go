ARG GO_VERSION=1.22.5
ARG RAILWAY_SERVICE_ID=algorhytm-go
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
LABEL org.opencontainers.image.source=https://github.com/rhymbic/algorhythm-go
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

FROM alpine:latest AS final

LABEL org.opencontainers.image.source=https://github.com/rhymbic/algorhythm-go

RUN apk --update add \
    ca-certificates \
    tzdata \
    && \
    update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

WORKDIR /app

COPY --from=build /bin/server /app/
COPY ./language/files ./language/files

ENTRYPOINT [ "/app/server" ]
