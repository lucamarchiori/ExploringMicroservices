FROM alpine:3.19.1
LABEL maintainer="Luca Marchiori" \
      description="Dockerfile for the API microservice"
WORKDIR /app
COPY /MicroserviceSource/bin /app
EXPOSE 4000
ENTRYPOINT [ "./server" ]

