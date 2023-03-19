FROM golang:1.19 AS build-env
WORKDIR /app
ADD . /app
RUN cd /app && CGO_ENABLED=0 go build -o app ./cmd/companies

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/app /app
COPY --from=build-env /app/migrations /app/migrations

EXPOSE 8080
CMD ["./app"]