FROM FROM alpine:3.2
RUN apk add --update ca-certificates # Certificates for SSL
COPY backend-services .
COPY docker-entrypoint.sh .
EXPOSE 8000
CMD ./docker-entrypoint.sh
