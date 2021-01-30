ARG go_version=1.16-alpine
FROM golang:$go_version AS build

WORKDIR /app

COPY $PWD/go.mod /app
COPY $PWD/go.sum /app
RUN go mod download

COPY $PWD /app
RUN CGO_ENABLED=0 go build -o pokemon-api

FROM scratch AS run

WORKDIR /app
COPY --from=build /app/pokemon-api /usr/bin/pokemon-api
COPY --from=build /etc/ssl/certs /etc/ssl/certs
EXPOSE 5000
ENTRYPOINT ["/usr/bin/pokemon-api"]
