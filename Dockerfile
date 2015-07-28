FROM busybox:ubuntu-14.04
COPY backend-services .
COPY docker-entrypoint.sh .
EXPOSE 8000
CMD ./docker-entrypoint.sh
