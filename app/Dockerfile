FROM golang:1.20 as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -v -o /usr/local/bin/app ./...

FROM scratch

WORKDIR /app

ENV GIN_MODE=release
EXPOSE 8080
COPY --from=build /usr/local/bin/app /app/clusterpuppy
COPY assets /app/assets
ENTRYPOINT ["/app/clusterpuppy"]
