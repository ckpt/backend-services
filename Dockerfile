FROM alpine:3.2
RUN apk add --update ca-certificates # Certificates for SSL
COPY backend-services /opt/ckpt/backend-services
COPY docker-entrypoint.sh /opt/ckpt/docker-entrypoint.sh
EXPOSE 8000
WORKDIR /opt/ckpt
CMD ./docker-entrypoint.sh
