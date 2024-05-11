FROM golang:1.21 AS BUILDER
WORKDIR /app
COPY ./src ./src
WORKDIR /app/src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/server ./cmd/api/

FROM alpine:3.19.1
LABEL maintainer="Luca Marchiori"
WORKDIR /app
COPY --from=0 /app/src/bin .
EXPOSE 4000
ENTRYPOINT [ "./server" ]


# docker build -t europe-west3-docker.pkg.dev/exploringmicroservices/exploringmicroservices/notification-ms:latest . && docker push europe-west3-docker.pkg.dev/exploringmicroservices/exploringmicroservices/notification-ms:latest