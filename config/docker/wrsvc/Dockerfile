FROM golang:1.16-alpine AS build

WORKDIR /src

ENV CGO_ENABLED 0
ENV GO111MODULE on

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build \
    -trimpath \
    -gcflags='-e -l' \
    -ldflags="-w -s -extldflags '-static'" \
    -o bin/wrsvc \
    ./cmd/wrsvc

FROM alpine:3.13 AS prod
COPY --from=build /src/bin/wrsvc /bin/wrsvc
ENTRYPOINT [ "wrsvc" ]
